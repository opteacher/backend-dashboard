import models from "./pages/models.vue"
import services from "./pages/services.vue"

export default [{
    path: "/",
    redirect: "/backend-dashboard/frontend/models"
}, {
    path: "/backend-dashboard/frontend/models",
    component: models
}, {
    path: "/backend-dashboard/frontend/services",
    component: services
}]