<template>
<div v-if="implInfo">
    <img v-holder="'img=100px200?auto=yes&theme=social'" class="card-img-top"/>
    <div class="card-body">
        <h5 class="card-title">{{implInfo.name}}</h5>
        <p class="card-text">{{implInfo.desc}}</p>
        <a :href="implInfo.homeUrl" class="btn btn-primary">官网链接</a>
    </div>
</div>
</template>

<script>
import backend from "../backend"

export default {
    props: {
        "implId": String
    },
    data() {
        return {
            implInfo: null
        }
    },
    async created() {
        let res = await backend.qryModSignById(this.implId)
        if (typeof res === "string") {
            this.$message.error(`查询模块信息时发生错误：${res}`)
        } else {
            this.implInfo = res
        }
    }
}
</script>