package v1

import (
	"net/http"

	"github.com/astronomerio/commander/provisioner"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("package", "v1")
)

// AirflowRouteHandler is a route handler for Airflow on Kubernetes
type AirflowRouteHandler struct {
	provisionHandler provisioner.Provisioner
}

// NewAirflowRouteHandler creates a new KubeAirflowRouteHandler
func NewAirflowRouteHandler(provHandler provisioner.Provisioner) *AirflowRouteHandler {
	return &AirflowRouteHandler{
		provisionHandler: provHandler,
	}
}

// Register registers the routes for KubeAirflowRouteHandler
func (h *AirflowRouteHandler) Register(router *gin.Engine) {
	airflowRouter := router.Group("/v1").Group("/airflow")
	{
		airflowRouter.GET(":organization_id/deployments", h.listDeployments)
		airflowRouter.POST(":organization_id/deployments", h.createDeployment)
		airflowRouter.PATCH(":organization_id/deployments/:deployment_id", h.patchDeployment)
		airflowRouter.DELETE(":organization_id/deployments/:deployment_id", h.deleteDeployment)
	}
}

func (h *AirflowRouteHandler) listDeployments(c *gin.Context) {
	logger := log.WithField("function", "listDeployments")
	logger.Debug("Entered listDeployments")

	organizationID := c.Param("organization_id")

	resp, listErr := h.provisionHandler.ListDeployments(organizationID)
	if listErr != nil {
		logger.Error(listErr)
		c.JSON(http.StatusInternalServerError, resp)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AirflowRouteHandler) createDeployment(c *gin.Context) {
	logger := log.WithField("function", "createDeployment")
	logger.Debug("Entered createDeployment")
	c.JSON(http.StatusOK, "")
}

func (h *AirflowRouteHandler) patchDeployment(c *gin.Context) {
	logger := log.WithField("function", "patchDeployment")
	logger.Debug("Entered patchDeployment")

	deploymentID := c.Param("deployment_id")

	patchReq := &provisioner.PatchDeploymentRequest{}
	if bindErr := c.BindJSON(patchReq); bindErr != nil {
		logger.Error(bindErr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, patchErr := h.provisionHandler.PatchDeployment(deploymentID, patchReq)
	if patchErr != nil {
		logger.Error(patchErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AirflowRouteHandler) deleteDeployment(c *gin.Context) {
	logger := log.WithField("function", "deleteDeployment")
	logger.Debug("Entered deleteDeployment")
	c.JSON(http.StatusOK, "")
}
