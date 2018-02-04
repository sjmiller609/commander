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

// DeploymentRouteHandler is a route handler for Airflow on Kubernetes
type DeploymentRouteHandler struct {
	provisionHandler provisioner.Provisioner
}

// NewDeploymentRouteHandler creates a new KubeAirflowRouteHandler
func NewDeploymentRouteHandler(provHandler provisioner.Provisioner) *DeploymentRouteHandler {
	return &DeploymentRouteHandler{
		provisionHandler: provHandler,
	}
}

// Register registers the routes for KubeAirflowRouteHandler
func (h *DeploymentRouteHandler) Register(router *gin.Engine) {
	depRouter := router.Group("/v1")
	{
		// depRouter.GET("/deployments", h.listDeployments)
		// depRouter.POST("/deployments", h.createDeployment)
		depRouter.PATCH("/deployments/:deploymentId/component/:componentId", h.patchDeployment)
		// depRouter.DELETE("deployments/:deployment_id", h.deleteDeployment)
	}

}

// func (h *DeploymentRouteHandler) listDeployments(c *gin.Context) {
// 	logger := log.WithField("function", "listDeployments")
// 	logger.Debug("Entered listDeployments")

// 	organizationID := c.Param("organization_id")
// 	resp, listErr := h.provisionHandler.ListDeployments(organizationID)
// 	if listErr != nil {
// 		logger.Error(listErr)
// 		c.JSON(http.StatusInternalServerError, resp)
// 	}

// 	c.JSON(http.StatusOK, resp)
// }

// func (h *DeploymentRouteHandler) createDeployment(c *gin.Context) {
// 	logger := log.WithField("function", "createDeployment")
// 	logger.Debug("Entered createDeployment")
// 	c.JSON(http.StatusOK, "")
// }

func (h *DeploymentRouteHandler) patchDeployment(c *gin.Context) {
	logger := log.WithField("function", "patchDeployment")
	logger.Debug("Entered patchDeployment")

	metadata := provisioner.DeploymentMetadata{
		DeploymentID: c.Param("deploymentId"),
		ComponentID:  c.Param("componentId"),
	}

	patchReq := provisioner.PatchDeploymentRequest{Metadata: metadata}
	if jsonErr := c.BindJSON(&patchReq); jsonErr != nil {
		logger.Error(jsonErr)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, patchErr := h.provisionHandler.PatchDeployment(&patchReq)
	if patchErr != nil {
		logger.Error(patchErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// func (h *DeploymentRouteHandler) deleteDeployment(c *gin.Context) {
// 	logger := log.WithField("function", "deleteDeployment")
// 	logger.Debug("Entered deleteDeployment")
// 	c.JSON(http.StatusOK, "")
// }
