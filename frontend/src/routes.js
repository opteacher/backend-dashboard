import home from "./pages/home.vue"

export default [{
    path: "/",
    redirect: "/backend-dashboard/frontend"
}, {
    path: "/backend-dashboard/frontend",
    component: home
}]