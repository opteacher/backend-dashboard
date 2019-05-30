<template>
<dashboard>
    <tool-bar @add-model="addModel" @add-relation="addRelation" :models="models"/>
    <div class="models-panel">
        <model-card v-for="model in models" :key="model.id" :model="model" @delete-model="deleteModel"/>
        <relation-link v-for="link in relations" :key="link.id" :relation="link"
                :model1="findModel(link.model1)"
                :model2="findModel(link.model2)"/>
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
        this.queryRelations()
    },
    watch: {
        relations() {
            // 注册为关联对象的观察者，观察模块的位置和尺寸的变化
            for (let relation of this.relations) {
                this.bindRelationToModel(relation.model1, relation)
                this.bindRelationToModel(relation.model2, relation)
            }
        }
    },
    methods: {
        async queryModels() {
            let res = await modelBkd.qry()
            if (typeof res === "string") {
                this.$message(`查询模块失败：${res}`)
            } else {
                this.models = res.data.data || []
            }
        },
        async queryRelations() {
            let res = await relationBkd.qry()
            if (typeof res === "string") {
                this.$message(`查询关联失败：${res}`)
            } else {
                this.relations = res.data.data || []
            }
        },
        async addModel(model) {
            let res = await modelBkd.add(model)
            if (typeof res === "string") {
                this.$message(`创建模块失败：${res}`)
            } else {
                model.id = res.data.data[0].id
                this.models.push(model)
            }
        },
        async deleteModel(modelID) {
            let res = await modelBkd.del(modelID)
            if (typeof res === "string") {
                this.$message(`删除模块失败：${res}`)
            } else {
                this.models.pop(ele => ele.id === modelID)
            }
        },
        async addRelation(relation) {
            let res = await relationBkd.add(relation)
            if (typeof res === "string") {
                this.$message(`创建关联失败：${res}`)
            } else {
                relation.id = res.data.data[0].id
                this.relations.push(relation)
            }
        },
        findModel(id) {
            return this.models.find(ele => ele.id === id)
        },
        bindRelationToModel(modelID, relation) {
            let model = this.findModel(modelID)
            model.observers = [relation]
            model.notifyUpdate = function() {
                for (let obs of this.observers) {
                    obs.onModelChanged(this.x, this.y, this.width, this.height)
                }
            }
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

