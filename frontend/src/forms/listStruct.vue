<template>
<el-collapse v-model="activeStruct" accordion>
    <el-collapse-item v-for="struct in structs" :key="struct.name" :name="struct.name">
        <template slot="title">
            <h6 class="mb-0">{{struct.name}}</h6>
        </template>
        <ul class="list-group">
            <li v-for="prop in struct.props" :key="prop.name" class="list-group-item d-flex justify-content-between align-items-center">
                {{prop.name}}<span class="badge badge-primary badge-pill">{{prop.type}}</span>
            </li>
        </ul>
        <button v-if="!baseNames.includes(struct.name)" type="button" class="btn btn-danger mt-10 w-100" @click="hdlDelStruct(struct.name)">删除</button>
    </el-collapse-item>
</el-collapse>
</template>

<script>
import backend from "../backend"
export default {
    props: {
        "showFlag": Boolean
    },
    data() {
        return {
            baseNames: [],
            structs: [],
            activeStruct: ""
        }
    },
    async created() {
        await this.refresh()
    },
    watch: {
        async showFlag() {
            if (this.showFlag) {
                await this.refresh()
            }
        }
    },
    methods: {
        async refresh() {
            let res = await backend.qryAllStructs()
            if (typeof res === "string") {
                this.$message.error(`查询所有结构时发生错误：${res}`)
            } else {
                this.structs = res.models ? res.models.map(struct => _.pick(struct, [
                    "name", "props"
                ])) : []
            }
            res = await backend.qryAllBaseStructsName()
            if (typeof res === "string") {
                this.$message.error(`查询基本结构名时发生错误：${res}`)
            } else {
                this.baseNames = res.names || []
            }
        },
        async hdlDelStruct(name) {
            this.$alert("确定删除结构？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delModel(name)
                    if (typeof res === "string") {
                        this.$message.error(`删除结构时发生错误：${res}`)
                    } else {
                        await this.refresh()
                    }
                    this.$message({
                        type: "info",
                        message: `结构（${name}）删除成功！`
                    })
                }
            })
        }
    }
}
</script>
