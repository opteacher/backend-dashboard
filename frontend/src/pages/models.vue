<template>
<dashboard>
    <tool-bar @add-model="addModel" @add-link="addLink" @add-struct="addStruct" :models="models"/>
    <div id="pnlModels" class="w-100 h-100" style="position: absolute">
        <div style="position: absolute; width: 0; height: 0">
            <model-card v-for="model in models" :key="model.name" :model="model"
                @delete-model="deleteModel"
                @update="updateLinksByModel"/>
        </div>
    </div>
    <svg id="pnlGraphs" class="w-100 h-100" style="position: absolute; z-index: -100"/>
</dashboard>
</template>

<script>
import utils from "../utils"
import dashboard from "../layouts/dashboard"
import toolBar from "../components/toolBar"
import modelCard from "../components/modelCard"
import backend from "../backend"

export default {
    components: {
        "dashboard": dashboard,
        "tool-bar": toolBar,
        "model-card": modelCard,
    },
    data() {
        return {
            models: [],
            links: []
        }
    },
    async created() {
        await this.queryModels()
        await this.queryLinks()
    },
    watch: {
        links() {
            let self = this
            d3.select("#pnlGraphs")
                .html("")
                .selectAll("g")
                .data(this.links)
                .join("g")
                .append("line")
                .attr("name", link => `link_${link.symbol}`)
                .attr("stroke-width", 1)
                .attr("stroke", "black")
                .each(function (link) {
                    self.updateLink(link, this)
                })
        }
    },
    methods: {
        updateLinksByModel(mname) {
            for (let link of this.links) {
                if (mname === link.mname1 || mname === link.mname2) {
                    this.updateLink(link)
                }
            }
        },
        updateLink(link, lnkInPnl) {
            let model1 = this.models.find(m => m.name === link.mname1)
            let model2 = this.models.find(m => m.name === link.mname2)
            link.x1 = model1.x + (model1.width>>1)
            link.y1 = model1.y + (model1.height>>1)
            link.x2 = model2.x + (model2.width>>1)
            link.y2 = model2.y + (model2.height>>1)
            if (!lnkInPnl) {
                lnkInPnl = `[name="link_${link.symbol}"]`
            }
            d3.select(lnkInPnl)
                .attr("x1", link => link.x1)
                .attr("y1", link => link.y1)
                .attr("x2", link => link.x2)
                .attr("y2", link => link.y2)
        },
        async queryModels() {
            let res = await backend.qryAllModels()
            if (typeof res === "string") {
                this.$message.error(`查询模块失败：${res}`)
            } else {
                this.models = res.models || []
            }
        },
        async queryLinks() {
            let res = await backend.qryAllLinks()
            if (typeof res === "string") {
                this.$message.error(`查询关联失败：${res}`)
            } else {
                this.links = res.links || []
            }
        },
        async addModel(model) {
            let res = await backend.addModel(model)
            if (typeof res === "string") {
                this.$message.error(`创建模块失败：${res}`)
            } else {
                this.models.push(model)
            }
        },
        async deleteModel(mname) {
            let res = await backend.delModel(mname)
            if (typeof res === "string") {
                this.$message.error(`删除模块失败：${res}`)
            } else {
                this.models.pop(ele => ele.name === mname)
                await this.queryLinks()
            }
        },
        async addLink(link) {
            link.symbol = `${link.mname1}-${link.mname2}`.toLowerCase()
            let res = await backend.addLink(link)
            if (typeof res === "string") {
                this.$message.error(`创建关联失败：${res}`)
            } else {
                link.id = res.id
                this.links.push(link)
            }
        },
        addStruct(struct) {
            struct = _.pick(struct, [
                "name", "props",
                "x", "y",
                "width", "height"
            ])
            console.log(struct)
        }
    },
}
</script>

<style lang="scss">
.models-panel {
    position: relative;
    height: 100%;
}
</style>

