import Vue from "vue"
import VueRouter from "vue-router"
import VueAxios from "vue-axios"
import axios from "axios"
import routes from "./routes"
import ElementUI from "element-ui"
import "./styles/common.css"
import VueHolder from 'vue-holderjs'

Vue.use(VueHolder)
Vue.use(VueRouter)
Vue.use(VueAxios, axios)
Vue.use(ElementUI)
Vue.router = new VueRouter({ routes })

new Vue({
    el: "#app",
	template: `
        <div class="h-100">
            <router-view></router-view>
        </div>`,
	router: Vue.router
})

export default routes.map(route => route.path)