<template>
<div class="card" :name="`flow_${flow.index}`" :style="`
    width:${flowWidth}px;margin-bottom:${flow.isLast ? marginTB : 0}px
`">
    <div class="card-header">
        #{{flow.index}}. <el-tag>{{flow.operKey}}</el-tag>
        <button class="close" type="button" data-dismiss="alert" aria-label="Close" @click="hdlDelete">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
    <div class="row">
        <div class="col pr-0">
            <ul class="list-group list-group-flush h-100">
                <a class="list-group-item list-group-item-primary list-group-item-action api-params" href="#" v-for="(content, pholder) in flow.inputs" :key="pholder">
                    {{pholder}}
                    <i class="el-icon-arrow-right"/>
                </a>
            </ul>
        </div>
        <div class="col-6 card-body text-center desc-panel" @click="showOperDetail">{{flow.desc}}</div>
        <div class="col pl-0">
            <div class="list-group list-group-flush h-100">
                <a class="list-group-item list-group-item-success list-group-item-action api-params text-right" href="#" v-for="output in flow.outputs" :key="output">
                    {{output}}
                    <i class="el-icon-arrow-right"/>
                </a>
            </div>
        </div>
    </div>
</div>
</template>

<script>
export default {
    props: {
        "flow": Object
    },
    data() {
        return {
            marginTB: 50,
            flowWidth: 600,
            flowSpan: 300,
        }
    },
    created() {
        let panelWidth = parseInt(document.getElementById("pnlFlows").getBoundingClientRect().width)
        this.flow.x = (panelWidth>>1) - (this.flowWidth>>1)
        this.flow.y = this.marginTB + this.flowSpan * this.flow.index
    },
    mounted() {
        d3.select(`[name="flow_${this.flow.index}"]`)
            .style("left", `${this.flow.x}px`)
            .style("top", `${this.flow.y}px`)
    },
    methods: {
        showOperDetail() {
            this.$emit("show-detail", this.flow)
        },
        hdlDelete() {
            // TODO: 删除步骤
        }
    }
}
</script>

