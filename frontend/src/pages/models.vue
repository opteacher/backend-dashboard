<template>
<dashboard>
    <tool-bar @add-model="addModel" @add-link="addLink" :models="models"/>
    <div id="pnlModels" class="w-100 h-100" style="position: absolute"></div>
    <svg id="pnlGraphs" class="w-100 h-100" style="position: absolute; z-index: -100" />
</dashboard>
</template>

<script>
import utils from "../utils"
import dashboard from "../layouts/dashboard"
import toolBar from "../components/toolBar"
import modelBkd from "../async/model"
import linkBkd from "../async/link"

export default {
        components: {
            "dashboard": dashboard,
            "tool-bar": toolBar,
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
            models() {
                let self = this
                let mdlPanel = d3.select("#pnlModels")
                    .html("")
                    .selectAll("div")
                    .data(this.models)
                    .join("div")
                    .style("position", "absolute")
                    .style("width", 0)
                    .style("height", 0)
                let mdlCard = mdlPanel.append("div")
                    .attr("class", "card")
                    .attr("name", model => `model_${model.name}`)
                    .style("left", model => `${model.x}px`)
                    .style("top", model => `${model.y}px`)
                    .style("width", model => `${model.width}px`)
                    .style("height", model => `${model.height}px`)
                    .style("cursor", "pointer")
                    .call(d3.drag().on("drag", function (tgt) {
                        tgt.x = d3.event.x >= 0 ? d3.event.x : 0
                        tgt.y = d3.event.y >= 0 ? d3.event.y : 0
                        d3.select(this)
                            .style("left", `${tgt.x}px`)
                            .style("top", `${tgt.y}px`)
                        for (let link of self.links) {
                            if (tgt.name === link.modelName1
                            || tgt.name === link.modelName2) {
                                self.updateLink(link)
                            }
                        }
                    }))
                mdlCard.append("div")
                    .attr("class", "card-header")
                    .text(model => model.name)
                    .append("button")
                    .attr("type", "button")
                    .attr("class", "close")
                    .attr("data-dismiss", "alert")
                    .attr("aria-label", "Close")
                    .on("click", eve => {
                        this.$alert("确定删除模块？", "提示", {
                            confirmButtonText: "确定",
                            callback: async action => {
                                if (action !== "confirm") {
                                    return
                                }
                                await this.deleteModel(eve.name)
                                this.$message({
                                    type: "info",
                                    message: `模块（${eve.name}）删除成功！`
                                })
                            }
                        })
                    })
                    .append("span")
                    .attr("aria-hidden", "true")
                    .html("&times;")
                mdlCard.append("div")
                    .attr("class", "card-body")
                    .append("ul")
                    .attr("class", "list-group list-group-flush")
                    .html(model => model.props.map(prop => `
                        <li class="list-group-item">
                            ${prop.name}
                            <span class="float-right">${prop.type}</span>
                        </li>
                    `).join(""))
                mdlCard.append("div")
                    .attr("class", "card-footer")
                    .html(model => [
                        {m:"POST", c:"success"},
                        {m:"DELETE", c:"danger"},
                        {m:"PUT", c:"warning"},
                        {m:"GET", c:"primary"},
                        {m:"ALL", c:"info"}
                    ].map(method => {
                        if (model.methods.includes(method.m)) {
                            return `<span class="badge badge-${method.c}">${method.m}</span> `
                        } else {
                            return `<span class="badge badge-secondary">${method.m}</span> `
                        }
                    }).join(""))
                let rszIcon = mdlCard.append("svg")
                    .attr("width", 18)
                    .attr("height", 18)
                    .style("position", "absolute")
                    .style("bottom", 0)
                    .style("right", 0)
                    .style("cursor", "nwse-resize")
                    .call(d3.drag().on("drag", function (tgt) {
                        let mouseLoc = d3.mouse(document.getElementById("pnlModels"))
                        d3.select(`[name="model_${tgt.name}"]`)
                            .style("width", `${tgt.width = mouseLoc[0] - tgt.x}px`)
                            .style("height", `${tgt.height = mouseLoc[1] - tgt.y}px`)
                        for (let link of self.links) {
                            if (tgt.name === link.modelName1
                            || tgt.name === link.modelName2) {
                                self.updateLink(link)
                            }
                        }
                    }))
                rszIcon.append("line")
                    .attr("x1", 16).attr("y1", 2)
                    .attr("x2", 2).attr("y2", 16)
                    .attr("stroke", "#7c7c7c")
                    .attr("stroke-width", 3)
                rszIcon.append("line")
                    .attr("x1", 16).attr("y1", 11)
                    .attr("x2", 11).attr("y2", 16)
                    .attr("stroke", "#7c7c7c")
                    .attr("stroke-width", 3)
            },
            links() {
                let self = this
                d3.select("#pnlGraphs")
                    .html("")
                    .selectAll("g")
                    .data(this.links)
                    .join("g")
                    .append("line")
                    .each(function (link) {
                        d3.select(this)
                            .attr("name", link => `link_${link.symbol}`)
                            .attr("stroke-width", 1)
                            .attr("stroke", "black")
                        self.updateLink(link, this)
                    })
            }
        },
        methods: {
            updateLink(link, lnkInPnl) {
                let model1 = this.models.find(m => m.name === link.modelName1)
                let model2 = this.models.find(m => m.name === link.modelName2)
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
                let res = await modelBkd.qry()
                if (typeof res === "string") {
                    this.$message.error(`查询模块失败：${res}`)
                } else {
                    this.models = (res.data.data && res.data.data.models) || []
                }
            },
            async queryLinks() {
                let res = await linkBkd.qry()
                if (typeof res === "string") {
                    this.$message.error(`查询关联失败：${res}`)
                } else {
                    this.links = res.data.data.links || []
                }
            },
            async addModel(model) {
                let res = await modelBkd.add(model)
                if (typeof res === "string") {
                    this.$message.error(`创建模块失败：${res}`)
                } else {
                    this.models.push(model)
                }
            },
            async deleteModel(mname) {
                let res = await modelBkd.del(mname)
                if (typeof res === "string") {
                    this.$message.error(`删除模块失败：${res}`)
                } else {
                    this.models.pop(ele => ele.name === mname)
                }
            },
            async addLink(link) {
                link.symbol = `${link.modelName1}-${link.modelName2}`.toLowerCase()
                let res = await linkBkd.add(link)
                if (typeof res === "string") {
                    this.$message.error(`创建关联失败：${res}`)
                } else {
                    link.id = res.data.data.id
                    this.links.push(link)
                }
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

