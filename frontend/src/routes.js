import models from "./pages/models.vue"
import services from "./pages/services.vue"
import steps from "./pages/steps"
import dao from "./pages/dao"
import axports from "./pages/exports"

export default [{
    path: "/",
    redirect: "/backend-dashboard/frontend/models"
}, {
    path: "/backend-dashboard/frontend/models",
    component: models
}, {
    path: "/backend-dashboard/frontend/services",
    component: services
}, {
    path: "/backend-dashboard/frontend/steps",
    component: steps
}, {
    path: "/backend-dashboard/frontend/dao",
    component: dao
}, {
    path: "/backend-dashboard/frontend/exports",
    component: axports
}]