<template>
<dashboard>
    <tool-bar @add-model="addModel" @add-relation="addRelation" :models="models"/>
    <div class="models-panel">
        <model-card v-for="model in models" :key="model.id" :model="model" @delete-model="deleteModel"/>
        <relation-link v-for="link in relations" :key="link.id" :model="link"
                :begModel="findModel(link.selBegMdl)"
                :endModel="findModel(link.selEndMdl)"/>
    </div>
</dashboard>
</template>

<script>
import utils from "../utils"
import modelBkd from "../async/model"
import relationBkd from "../async/relation"

import dashboard from "../layouts/dashboard"
import toolBar from "../components/toolBar"
import modelCard from "../components/modelCard"
import relationLink from "../components/relationLink"

export default {
    data() { return {
        movingModel: null,
        models: [],
        relations: []
    }},
    created() {
        this.queryModels()
    },
    methods: {
        async queryModels() {
            let res = await modelBkd.get()
            if (typeof res === "string") {
                this.$message(`查询模块失败：${res}`)
            } else {
                this.models = res.data.data || []
            }
        },
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
        },
        async addRelation(relation) {
            let res = await relationBkd.post(relation)
            if (typeof res === "string") {
                this.$message(`创建关联失败：${res}`)
            } else {
                this.relations.push(res.data.data)
            }
        },
        findModel(id) {
            return this.models.find(ele => ele.id === id)
        }
    },
    components: {
        "dashboard": dashboard,
        "model-card": modelCard,
        "tool-bar": toolBar,
        "relation-link": relationLink
    }
}
</script>

<style lang="scss">
.models-panel {
    position: relative;
    height: 100%;
}
</style>

