<template>
<dashboard>
    <tool-bar @add-model="addModel"/>
    <div class="models-panel">
        <model-card v-for="model in models" :key="model.id" :model="model" @delete-model="deleteModel"/>
    </div>
</dashboard>
</template>

<script>
import utils from "../utils"

import dashboard from "../layouts/dashboard"
import toolBar from "../components/toolBar"
import modelCard from "../components/modelCard"
import modelBkd from "../async/model"

export default {
    data() {
        return {
            movingModel: null,
            models: []
        }
    },
    methods: {
        async addModel(model) {
            let res = await modelBkd.post(model)
            if (typeof res === "string") {
                this.$message(`创建模块失败：${res}`)
            } else {
                model.id = res.data.data[0].id
                this.models.push(model)
            }
        },
        async deleteModel(modelID) {
            let res = await modelBkd.delete(modelID)
            if (typeof res === "string") {
                this.$message(`删除模块失败：${res}`)
            } else {
                this.models.pop(ele => ele.id === modelID)
            }
        }
    },
    components: {
        "dashboard": dashboard,
        "model-card": modelCard,
        "tool-bar": toolBar,
    }
}
</script>

<style lang="scss">
.models-panel {
    position: relative;
    height: 100%;
}
</style>

