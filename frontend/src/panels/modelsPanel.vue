<template>
    <div class="w-100 h-100">
        <tool-bar @add-model="addModel" @add-link="addLink" :models="models"/>
        <div id="pnlModels" class="w-100 h-100"></div>
    </div>
</template>

<script>
    import toolBar from "../components/toolBar"
    import modelBkd from "../async/model"
    import relationBkd from "../async/relation"

    export default {
        components: {
            "tool-bar": toolBar,
        },
        data() {
            return {
                inited: false,
                models: [],
                links: []
            }
        },
        created() {
            this.queryModels()
            // this.queryLinks()
        },
        watch: {
            models() {
                d3.select("#pnlModels").html("")
                let mdlPanel = d3.select("#pnlModels")
                    .selectAll("div")
                    .data(this.models)
                    .join("div")
                    .style("position", "absolute")
                let mdlCard = mdlPanel.append("div")
                    .attr("class", "card")
                    .attr("name", model => `model_${model.name}`)
                    .style("left", model => `${model.x}px`)
                    .style("top", model => `${model.y}px`)
                    .style("width", model => `${model.width}px`)
                    .style("height", model => `${model.height}px`)
                    .style("cursor", "pointer")
                    .call(d3.drag().on("drag", function (tgt) {
                        d3.select(this)
                            .style("left", `${tgt.x = d3.event.x}px`)
                            .style("top", `${tgt.y = d3.event.y}px`)
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
            }
        },
        methods: {
            async queryModels() {
                let res = await modelBkd.qry()
                if (typeof res === "string") {
                    this.$message(`查询模块失败：${res}`)
                } else {
                    this.models = (res.data.data && res.data.data.models) || []
                }
            },
            async queryLinks() {
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
                    this.models.push(model)
                }
            },
            async deleteModel(mname) {
                let res = await modelBkd.del(mname)
                if (typeof res === "string") {
                    this.$message(`删除模块失败：${res}`)
                } else {
                    this.models.pop(ele => ele.name === mname)
                }
            },
            async addLink(relation) {
                let res = await relationBkd.add(relation)
                if (typeof res === "string") {
                    this.$message(`创建关联失败：${res}`)
                } else {
                    relation.id = res.data.data[0].id
                    this.relations.push(relation)
                }
            }
        },
    }
</script>
