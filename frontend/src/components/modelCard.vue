<template>
<el-card class="box-card model-card" :style="`left: ${model.x}px; top: ${model.y}px; width: ${model.width}px; height: ${model.height}px`">
    <div slot="header" class="clearfix" v-drag="model.id">
        <span class="card-name">{{model.name}}</span>
        <el-link class="float-right" @click="$emit('delete-model', model.id)">
            <i class="el-icon-close"/>
        </el-link>
    </div>
    <div class="text item">
        <div class="list-group list-group-flush">
            <button type="button" class="list-group-item list-group-item-action list-group-item-light" v-for="prop in model.props" :key="prop.name">{{prop.name}}</button>
        </div>
    </div>
    <el-link class="resize-button" type="primary" v-resize="model.id">
        <i class="iconfont icon--Resize-Four-Direc"/>
    </el-link>
</el-card>
</template>

<script>
import $ from "jquery"

import modelBkd from "../async/model"

export default {
    props: {
        "model": Object
    },
    created() {
        console.log(this.model)
    },
    directives: {
        drag: { bind(el, binding) {
            let id = binding.value
            el.onmousedown = me => {
                let card = $(el).closest(".model-card")[0]
                let left = Number(card.style.left.slice(0, -2))
                let top = Number(card.style.top.slice(0, -2))
                let downClientX = me.clientX
                let downClientY = me.clientY

                document.onmousemove = e => {
                    let l = left + (e.clientX - downClientX)
                    let t = top + (e.clientY - downClientY)
                    t = (t <= 0) ? 0 : t
                    l = (l <= 0) ? 0 : l
                    card.style.left = `${l}px`
                    card.style.top = `${t}px`
                }
                document.onmouseup = async e => {
                    document.onmousemove = null
                    document.onmouseup = null
                    let res = await modelBkd.put({
                        id,
                        x: Number(card.style.left.slice(0, -2)),
                        y: Number(card.style.top.slice(0, -2))
                    })
                    if (typeof res === "string") {
                        console.log(`更新模块失败：${res}`)
                    }
                }
            }
        }},
        resize: { bind(el, binding) {
            let id = binding.value
            el.onmousedown = me => {
                let card = $(el).closest(".model-card")[0]
                let width = Number(card.style.width.slice(0, -2))
                let height = Number(card.style.height.slice(0, -2))
                let downClientX = me.clientX
                let downClientY = me.clientY

                document.onmousemove = e => {
                    let w = width +  e.clientX - downClientX
                    let h = height + e.clientY - downClientY
                    h = (h <= 61 + 24 + 40) ? 61 + 24 + 40 : h
                    card.style.width = `${w}px`
                    card.style.height = `${h}px`
                }
                document.onmouseup = async e => {
                    document.onmousemove = null
                    document.onmouseup = null
                    let res = await modelBkd.put({
                        id,
                        width: Number(card.style.width.slice(0, -2)),
                        height: Number(card.style.height.slice(0, -2))
                    })
                    if (typeof res === "string") {
                        console.log(`更新模块失败：${res}`)
                    }
                }
            }
        }}
    }
}
</script>


<style lang="scss">
.model-card {
    position: absolute;
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;

    .el-card__header {
        cursor: pointer;
    }
}
.resize-button {
    position: absolute;
    right: 20px;
    bottom: 20px;
}
</style>
