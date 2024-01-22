package delivery

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/auth"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"

	commonHttp "github.com/yarikTri/web-transport-cards/internal/common/http"
)

const SERVICE_TICKET_NAME = "X-Service-Ticket"
const SERVICE_TICKET_CORRECT_VALUE = "Q5%7&fG*"

type Handler struct {
	ticketServices ticket.Usecase
	authServices   auth.Usecase
	logger         logger.Logger
}

func NewHandler(tu ticket.Usecase, au auth.Usecase, l logger.Logger) *Handler {
	return &Handler{
		ticketServices: tu,
		authServices:   au,
		logger:         l,
	}
}

// @Summary		Get ticket
// @Tags		Tickets
// @Description	Get ticket by ID
// @Produce     json
// @Param		ticketID path int true 							"Ticket ID"
// @Success		200			{object}	models.TicketTransfer	"Got ticket"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/tickets/{ticketID} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	isService := c.GetHeader("X-SERVICE") == "true"

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil && !isService {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusUnauthorized, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil && !isService {
		h.logger.Infof("User not found")
		c.JSON(http.StatusInternalServerError, "User not found")
		return
	}

	isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID))

	ticket, err := h.ticketServices.GetByID(int(id), int(user.ID))
	if err != nil {
		h.logger.Errorf("Error while getting ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if !isModerator && !isService && ticket.CreatorID != int(user.ID) {
		h.logger.Error("Forbidden to get ticket")
		c.JSON(http.StatusForbidden, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		List tickets
// @Tags		Tickets
// @Description	Get all not draft tickets
// @Produce     json
// @Success		200			{object}	[]models.TicketTransfer	"Got tickets"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets [get]
func (h *Handler) List(c *gin.Context) {
	formTimeStartQuery, _ := strconv.ParseInt(c.Query("formTimeStart"), 10, 32)
	formTimeEndQuery, _ := strconv.ParseInt(c.Query("formTimeEnd"), 10, 32)

	stateQuery := c.Query("state")

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusForbidden, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID))

	tickets, err := h.ticketServices.List()
	if err != nil {
		h.logger.Errorf("Error while listing tickets: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	draft, _ := h.ticketServices.GetDraft(int(user.ID))

	tickets = append(tickets, draft)

	ticketsTransfers := make([]models.TicketTransfer, 0)
	for _, ticket := range tickets {
		formTime := ticket.FormTime.Unix()

		formTimeExpr := (formTimeStartQuery == 0 || formTime >= formTimeStartQuery) && (formTimeEndQuery == 0 || formTime <= formTimeEndQuery)
		stateExpr := stateQuery == "" || ticket.State == stateQuery
		userExpr := isModerator || ticket.CreatorID == int(user.ID)

		if formTimeExpr && stateExpr && userExpr {
			ticketsTransfers = append(ticketsTransfers, ticket.ToTransfer())
		}
	}

	c.JSON(http.StatusOK, ticketsTransfers)
}

// @Summary		Form ticket
// @Tags		Tickets
// @Description	Form ticket draft by ID
// @Produce     json
// @Param		ticketID path int true 							"Ticket ID"
// @Success		200			{object}	models.TicketTransfer	"Ticket formed"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets/draft/form [put]
func (h *Handler) FormDraft(c *gin.Context) {
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusBadRequest, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	ticket, err := h.ticketServices.FormDraft(int(user.ID))
	if err != nil {
		h.logger.Errorf("Error while forming draft ticket of user with id %d: %w", user.ID, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		End ticket
// @Tags		Tickets
// @Description	End ticket by ID
// @Produce     json
// @Param		ticketID path int true 							"Ticket ID"
// @Success		200			{object}	models.TicketTransfer	"Ticket ended"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets/{ticketID}/end [put]
func (h *Handler) EndByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusForbidden, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	if isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID)); !isModerator {
		h.logger.Infof("User is not a moderator")
		c.JSON(http.StatusForbidden, "User is not a moderator")
		return
	}

	ticket, err := h.ticketServices.EndByID(int(id), int(user.ID))
	if err != nil {
		h.logger.Errorf("Error while ending ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		Moderate ticket
// @Tags		Tickets
// @Description	Moderate formed ticket by ID
// @Accept 		json
// @Produce     json
// @Param		ticketID path int true 							"Ticket ID"
// @Param		req body UpdateStateRequest true 				"Ticket new state"
// @Success		200			{object}	models.TicketTransfer	"Ticket moderated"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets/{ticketID}/moderate [put]
func (h *Handler) ModerateByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusUnauthorized, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	if isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID)); !isModerator {
		h.logger.Infof("User is not a moderator")
		c.JSON(http.StatusForbidden, "User is not a moderator")
		return
	}

	var ticket models.Ticket

	var req UpdateStateRequest
	c.BindJSON(&req)
	switch req.NewState {
	case "approved":
		ticket, err = h.ticketServices.ApproveByID(int(id), int(user.ID))
		break
	case "rejected":
		ticket, err = h.ticketServices.RejectByID(int(id), int(user.ID))
		break
	default:
		h.logger.Infof("Invalid tickets's state '%s'", req.NewState)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket's state '%s'", req.NewState))
		return
	}

	if err != nil {
		h.logger.Errorf("Error while rejecting or approving ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		Delete ticket
// @Tags		Tickets
// @Description	Delete ticket by ID
// @Produce     json
// @Param		ticketID path int true 			"Ticket ID"
// @Success		200								"Ticket deleted"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/tickets/{ticketID} [delete]
func (h *Handler) DeleteByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	if err := h.ticketServices.DeleteByID(int(id)); err != nil {
		h.logger.Errorf("Error while getting ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Add route to ticket
// @Tags		Routes
// @Description	Add route to ticket draft by ID
// @Produce     json
// @Param		routeID path int true 							"Route ID"
// @Success		200			{object}	models.TicketTransfer	"Route added"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets/routes/{routeID} [post]
func (h *Handler) AddRoute(c *gin.Context) {
	routeID, err := strconv.ParseUint(c.Param("route_id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("route_id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusBadRequest, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	_, err = h.getOrCreateDraftTicket(int(user.ID))
	if err != nil {
		h.logger.Errorf("Error while getting draft ticket: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ticket, err := h.ticketServices.AddRoute(int(user.ID), int(routeID))
	if err != nil {
		h.logger.Errorf("Error while adding route with id %d to draft of user with id %d: %w", routeID, user.ID, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		Delete route from ticket
// @Tags		Routes
// @Description	Delete route from ticket draft
// @Produce     json
// @Param		ticketID path int true 							"Ticket ID"
// @Param		routeID path int true 							"Route ID"
// @Success		200			{object}	models.TicketTransfer	"Route deleted from ticket draft"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets/routes/{routeID} [delete]
func (h *Handler) DeleteRoute(c *gin.Context) {
	routeID, err := strconv.ParseUint(c.Param("route_id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("route_id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid route id '%s'", c.Param("id")))
		return
	}

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusBadRequest, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	_, err = h.getOrCreateDraftTicket(int(user.ID))
	if err != nil {
		h.logger.Errorf("Error while getting draft ticket: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ticket, err := h.ticketServices.DeleteRoute(int(user.ID), int(routeID))
	if err != nil {
		h.logger.Errorf("Error while deleting route with id %d from draft ticket of user with id %d: %w", routeID, user.ID, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		Get ticket draft
// @Tags		Tickets
// @Description	Get ticket draft
// @Produce     json
// @Success		200			{object}	models.TicketTransfer	"Ticket draft for current user"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/tickets/draft [get]
func (h *Handler) GetDraft(c *gin.Context) {
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusUnauthorized, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	ticket, err := h.ticketServices.GetDraft(int(user.ID))
	if err != nil {
		h.logger.Infof("Error while getting draft of user with id '%d'", user.ID)
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// @Summary		Delete ticket draft
// @Tags		Tickets
// @Description	Delete ticket draft
// @Success		200								"Draft deleted"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/tickets/draft [delete]
func (h *Handler) DeleteDraft(c *gin.Context) {
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusBadRequest, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	if err := h.ticketServices.DeleteDraft(int(user.ID)); err != nil {
		h.logger.Infof("Error while deleting draft of user with id '%d'", user.ID)
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Start updating ticket's write state
// @Tags		Tickets
// @Description	Start updating ticket's write state
// @Success		200								"Write state started to update"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/tickets/{ticket_id}/start_update_write_state [put]
func (h *Handler) StartWriting(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusBadRequest, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	if isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID)); !isModerator {
		h.logger.Infof("User is not a moderator")
		c.JSON(http.StatusForbidden, "User is not a moderator")
		return
	}

	t, err := h.ticketServices.GetByID(int(id), int(user.ID))
	if err != nil {
		h.logger.Infof("Ticket with id %d not found", id)
		c.JSON(http.StatusNotFound, "Ticket with such id not found")
		return
	}

	if t.State != models.APPROVED_STATE {
		h.logger.Infof("Invalid ticket's state to update write state")
		c.JSON(http.StatusForbidden, "Invalid ticket's state to update write state")
		return
	}

	postEndpoint := fmt.Sprintf("http://127.0.0.1:8000/write_ticket/%d", id)
	fmt.Println("Post endpoint: ", postEndpoint)

	resp, err := http.Post(postEndpoint, "application/json", bytes.NewBuffer(nil))
	if err != nil {
		fmt.Println("Resp: ", resp)
		fmt.Println("Err: ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(resp.StatusCode, err)
}

// @Summary		Update ticket's write state
// @Tags		Tickets
// @Description	Update ticket's write state
// @Success		200								"Write state updated"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/tickets/{ticket_id}/update_write_state [put]
func (h *Handler) UpdateWriteState(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	if isService := c.GetHeader(SERVICE_TICKET_NAME) == SERVICE_TICKET_CORRECT_VALUE; !isService {
		h.logger.Infof("Invalid service ticket")
		c.JSON(http.StatusForbidden, fmt.Sprintf("Invalid service ticket"))
		return
	}

	var req FinalizeWritingRequest
	c.BindJSON(&req)
	state := req.State
	if state != "success" && state != "fail" && state != "updating" {
		h.logger.Infof("Invalid state '%s'", req.State)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid state '%s'", req.State))
		return
	}

	updatedTicket, err := h.ticketServices.UpdateWriteState(int(id), req.State)
	if err != nil {
		h.logger.Infof("%w", err)
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, updatedTicket.ToTransfer())
}

func (h *Handler) getOrCreateDraftTicket(userID int) (int, error) {
	var ticketID int
	foundTicket, err := h.ticketServices.GetDraft(userID)
	if err != nil {
		// Черновик транспортной карты не найден - создаём
		fmt.Println("СОЗДАЁМ ЧЕРНОВИК")
		ticket, err := h.ticketServices.Create(
			models.Ticket{
				CreatorID: int(userID),
				State:     models.DRAFT_STATE,
			},
		)
		if err != nil {
			return 0, err
		}
		ticketID = int(ticket.ID)
	} else {
		ticketID = int(foundTicket.ID)
	}

	return ticketID, nil
}
