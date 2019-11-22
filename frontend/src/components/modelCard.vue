<template>
<div style="position:absolate;left:0;top:0">
    <div class="card" :name="`model_${model.name}`" :style="`
        cursor:pointer;left:${model.x}px;top:${model.y}px;width:${model.width}px;height:${model.height}px
    `">
        <div class="card-header">
            {{model.name}}
            <button class="close" type="button" data-dismiss="alert" aria-label="Close" @click="hdlDelete">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>
        <div class="card-body">
            <ul class="list-group list-group-flush">
                <li class="list-group-item" v-for="prop in model.props" :key="prop.name">
                    {{prop.name}}
                    <span class="float-right">{{prop.type}}</span>
                </li>
            </ul>
        </div>
        <div class="card-footer">
            <span v-for="method in [
                {m:'POST', c:'success'},
                {m:'DELETE', c:'danger'},
                {m:'PUT', c:'warning'},
                {m:'GET', c:'primary'},
                {m:'ALL', c:'info'}
            ]" :key="method.m" :class="`mr-1 badge badge-${model.methods.includes(method.m) ? method.c : 'secondary'}`">
                {{method.m}}
            </span>
            <svg :name="`model_${model.name}_resize`" width="18" height="18" style="position: absolute; bottom: 0; right: 0; cursor: nwse-resize">
                <line x1="16" y1="2" x2="2" y2="16" stroke="#7c7c7c" stroke-width="3"/>
                <line x1="16" y1="11" x2="11" y2="16" stroke="#7c7c7c" stroke-width="3"/>
            </svg>
        </div>
    </div>
</div>
</template>

<script>
import _ from "lodash"
import backend from "../backend"
import modelCard from "../components/modelCard"

export default {
    props: {
        "model": Object
    },
    mounted() {
        let self = this
        let model = this.model
        d3.select(`[name="model_${this.model.name}"]`)
            .call(d3.drag().on("start", function() {
                model.dx = d3.event.x - model.x
                model.dy = d3.event.y - model.y
            }).on("drag", function() {
                d3.select(this)
                    .style("left", `${model.x = d3.event.x >= 0 ? d3.event.x - model.dx : 0}px`)
                    .style("top", `${model.y = d3.event.y >= 0 ? d3.event.y - model.dy : 0}px`)
                self.$emit("update", model.name)
            }).on("end", async function() {
                let res = await backend.updModel(_.pick(model, ["name", "x", "y", "width", "height"]))
                if (typeof res === "string") {
                    self.$message.error(`更新模型位置时发生错误：${res}`)
                }
            }))
        d3.select(`[name="model_${this.model.name}_resize"]`)
            .call(d3.drag().on("drag", function () {
                let mouseLoc = d3.mouse(document.getElementById("pnlModels"))
                d3.select(`[name="model_${model.name}"]`)
                    .style("width", `${model.width = mouseLoc[0] - model.x}px`)
                    .style("height", `${model.height = mouseLoc[1] - model.y}px`)
                self.$emit("update", model.name)
            }).on("end", async function() {
                let res = await backend.updModel(_.pick(model, ["name", "x", "y", "width", "height"]))
                if (typeof res === "string") {
                    self.$message.error(`更新模型尺寸时发生错误：${res}`)
                }
            }))
    },
    methods: {
        hdlDelete() {
            this.$alert("确定删除模块？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    this.$emit("delete-model", this.model.name)
                    this.$message({
                        type: "info",
                        message: `模块（${this.model.name}）删除成功！`
                    })
                }
            })
        }
    }
}
</script>
