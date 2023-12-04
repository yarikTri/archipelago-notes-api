package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"
)

type Handler struct {
	services ticket.Usecase
	logger   logger.Logger
}

func NewHandler(tu ticket.Usecase, l logger.Logger) *Handler {
	return &Handler{
		services: tu,
		logger:   l,
	}
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	ticket, err := h.services.GetByID(int(id))
	if err != nil {
		h.logger.Errorf("Error while getting ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

func (h *Handler) List(c *gin.Context) {
	tickets, err := h.services.List()
	if err != nil {
		h.logger.Errorf("Error while listing tickets: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ticketsTransfers := make([]models.TicketTransfer, 0)
	for _, ticket := range tickets {
		ticketsTransfers = append(ticketsTransfers, ticket.ToTransfer())
	}

	c.JSON(http.StatusOK, ticketsTransfers)
}

func (h *Handler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ticket": "update"})
}

func (h *Handler) FormByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	ticket, err := h.services.FormByID(int(id))
	if err != nil {
		h.logger.Errorf("Error while forming ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

// reject/approve by moderator
func (h *Handler) ModerateByID(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"message": "You can't approve or reject ticket: no moderator's rights"})
	// return

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	var ticket models.Ticket

	var req UpdateStateRequest
	c.BindJSON(&req)
	switch req.NewState {
	case "approved":
		ticket, err = h.services.ApproveByID(int(id))
		break
	case "rejected":
		ticket, err = h.services.RejectByID(int(id))
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

func (h *Handler) DeleteByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	if err := h.services.DeleteByID(int(id)); err != nil {
		h.logger.Errorf("Error while getting ticket with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) AddRoute(c *gin.Context) {
	routeID, err := strconv.ParseUint(c.Param("route_id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("route_id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	var ticketID int
	foundTicket := h.services.GetTicketDraftByCreatorID(models.DEFAULT_CREATOR_ID)
	if foundTicket == nil {
		// Транспортная карта не найдена - создаём черновик
		ticket, err := h.services.Create(
			models.Ticket{
				CreatorID: models.DEFAULT_CREATOR_ID,
				State:     models.DRAFT_STATE,
			},
		)
		if err != nil {
			h.logger.Errorf("Error while creating draft ticket: %w", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		ticketID = int(ticket.ID)
	} else {
		ticketID = int(foundTicket.ID)
	}

	ticket, err := h.services.AddRoute(int(ticketID), int(routeID))
	if err != nil {
		h.logger.Errorf("Error while adding route with id %d to ticket with id %d: %w", ticketID, routeID, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}

func (h *Handler) DeleteRoute(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid ticket id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid ticket id '%s'", c.Param("id")))
		return
	}

	routeID, err := strconv.ParseUint(c.Param("route_id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("route_id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid route id '%s'", c.Param("id")))
		return
	}

	ticket, err := h.services.DeleteRoute(int(ticketID), int(routeID))
	if err != nil {
		h.logger.Errorf("Error while deleting route with id %d from ticket with id %d: %w", ticketID, routeID, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ticket.ToTransfer())
}
