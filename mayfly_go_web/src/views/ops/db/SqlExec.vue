<template>
    <div>
        <div class="toolbar">
            <el-row type="flex" justify="space-between">
                <el-col :span="24">
                    <project-env-select @changeProjectEnv="changeProjectEnv" @clear="clearDb">
                        <template #default>
                            <el-form-item label="数据库">
                                <el-select v-model="dbId" placeholder="请选择数据库" @change="changeDb" @clear="clearDb" clearable filterable>
                                    <el-option v-for="item in dbs" :key="item.id" :label="item.database" :value="item.id">
                                        <span style="float: left">{{ item.database }}</span>
                                        <span style="float: right; color: #8492a6; margin-left: 6px; font-size: 13px">{{
                                            `${item.name}  [${item.type}]`
                                        }}</span>
                                    </el-option>
                                </el-select>
                            </el-form-item>
                        </template>
                    </project-env-select>
                </el-col>
            </el-row>
        </div>

        <el-container style="border: 1px solid #eee; margin-top: 1px">
            <el-aside id="sqlcontent" width="65%" style="background-color: rgb(238, 241, 246)">
                <div class="toolbar">
                    <div class="fl">
                        <el-button @click="runSql" type="success" icon="el-icon-video-play" size="mini" plain>执行</el-button>

                        <el-button @click="formatSql" type="primary" icon="el-icon-magic-stick" size="mini" plain>格式化</el-button>

                        <el-button @click="saveSql" type="primary" icon="el-icon-document-add" size="mini" plain>保存</el-button>
                    </div>
                </div>
                <codemirror @beforeChange="onBeforeChange" class="codesql" ref="cmEditor" language="sql" v-model="sql" :options="cmOptions" />
            </el-aside>

            <el-container style="margin-left: 2px">
                <el-header style="text-align: left; height: 45px; font-size: 12px; padding: 0px">
                    <el-select v-model="tableName" placeholder="请选择表" @change="changeTable" filterable style="width: 99%">
                        <el-option
                            v-for="item in tableMetadata"
                            :key="item.tableName"
                            :label="item.tableName + (item.tableComment != '' ? `【${item.tableComment}】` : '')"
                            :value="item.tableName"
                        >
                        </el-option>
                    </el-select>
                </el-header>

                <el-main style="padding: 0px; overflow: hidden">
                    <el-table :data="columnMetadata" height="100%" size="mini">
                        <el-table-column prop="columnName" label="名称" show-overflow-tooltip> </el-table-column>
                        <el-table-column prop="columnComment" label="备注" show-overflow-tooltip> </el-table-column>
                        <el-table-column width="120" prop="columnType" label="类型" show-overflow-tooltip> </el-table-column>
                    </el-table>
                </el-main>
            </el-container>
        </el-container>

        <el-table style="margin-top: 1px" :data="execRes.data" size="mini" max-height="300" :empty-text="execRes.emptyResText" stripe border>
            <el-table-column
                min-width="100"
                :width="flexColumnWidth(item, execRes.data)"
                align="center"
                v-for="item in execRes.tableColumn"
                :key="item"
                :prop="item"
                :label="item"
                show-overflow-tooltip
            >
            </el-table-column>
        </el-table>
    </div>
</template>

<script lang="ts">
import { toRefs, reactive, computed, defineComponent, ref } from 'vue';
import { dbApi } from './api';

import 'codemirror/theme/ambiance.css';
import 'codemirror/addon/hint/show-hint.css';
// import base style
import 'codemirror/lib/codemirror.css';
// 引入主题后还需要在 options 中指定主题才会生效
import 'codemirror/theme/base16-light.css';

import 'codemirror/addon/selection/active-line';
import { codemirror } from '@/components/codemirror';
// import 'codemirror/mode/sql/sql.js';
import 'codemirror/addon/hint/show-hint.js';
import 'codemirror/addon/hint/sql-hint.js';

import sqlFormatter from 'sql-formatter';
import { notNull, notEmpty } from '@/common/assert';
import { ElMessage } from 'element-plus';
import ProjectEnvSelect from '../component/ProjectEnvSelect.vue';

export default defineComponent({
    name: 'SqlExec',
    components: {
        codemirror,
        ProjectEnvSelect,
    },
    setup() {
        const cmEditor: any = ref(null);
        const state = reactive({
            dbs: [],
            tables: [],
            dbId: null,
            tableName: '',
            tableMetadata: [],
            columnMetadata: [],
            sql: '',
            sqlTabs: {
                tabs: [] as any,
                active: '',
                index: 1,
            },
            execRes: {
                tableColumn: [],
                data: [],
                emptyResText: '没有数据',
            },
            params: {
                pageNum: 1,
                pageSize: 10,
                envId: null,
            },
            cmOptions: {
                tabSize: 4,
                mode: 'text/x-sql',
                lineNumbers: true,
                line: true,
                indentWithTabs: true,
                smartIndent: true,
                matchBrackets: true,
                theme: 'base16-light',
                autofocus: true,
                extraKeys: { Tab: 'autocomplete' }, // 自定义快捷键
                hintOptions: {
                    completeSingle: false,
                    // 自定义提示选项
                    tables: {},
                },
                // more CodeMirror options...
            },
        });

        const codemirror: any = computed(() => {
            return cmEditor.value.coder;
        });

        /**
         * 项目及环境更改后的回调事件
         */
        const changeProjectEnv = (projectId: any, envId: any) => {
            state.dbs = [];
            state.dbId = null;
            clearDb();
            if (envId != null) {
                state.params.envId = envId;
                search();
            }
        };

        /**
         * 输入字符给提示
         */
        const inputRead = (instance: any, changeObj: any) => {
            if (/^[a-zA-Z]/.test(changeObj.text[0])) {
                showHint();
            }
        };

        const onBeforeChange = (instance: any, changeObj: any) => {
            var text = changeObj.text[0];
            // 将sql提示去除
            changeObj.text[0] = text.split('  ')[0];
        };

        /**
         * 执行sql
         */
        const runSql = async () => {
            notNull(state.dbId, '请先选择数据库');
            // 没有选中的文本，则为全部文本
            let sql = getSql();
            notNull(sql, '内容不能为空');

            state.execRes.tableColumn = [];
            state.execRes.data = [];
            state.execRes.emptyResText = '查询中...';

            const res = await dbApi.sqlExec.request({
                id: state.dbId,
                sql: sql,
            });
            state.execRes.emptyResText = '没有数据';
            state.execRes.tableColumn = res.colNames;
            state.execRes.data = res.res;
        };

        const flexColumnWidth = (str: any, tableData: any, flag = 'equal') => {
            // str为该列的字段名(传字符串);tableData为该表格的数据源(传变量);
            // flag为可选值，可不传该参数,传参时可选'max'或'equal',默认为'max'
            // flag为'max'则设置列宽适配该列中最长的内容,flag为'equal'则设置列宽适配该列中第一行内容的长度。
            str = str + '';
            let columnContent = '';
            if (!tableData || !tableData.length || tableData.length === 0 || tableData === undefined) {
                return;
            }
            if (!str || !str.length || str.length === 0 || str === undefined) {
                return;
            }
            if (flag === 'equal') {
                // 获取该列中第一个不为空的数据(内容)
                for (let i = 0; i < tableData.length; i++) {
                    if (tableData[i][str].length > 0) {
                        columnContent = tableData[i][str];
                        break;
                    }
                }
            } else {
                // 获取该列中最长的数据(内容)
                let index = 0;
                for (let i = 0; i < tableData.length; i++) {
                    if (tableData[i][str] === null) {
                        return;
                    }
                    const now_temp = tableData[i][str] + '';
                    const max_temp = tableData[index][str] + '';
                    if (now_temp.length > max_temp.length) {
                        index = i;
                    }
                }
                columnContent = tableData[index][str];
            }
            // 以下分配的单位长度可根据实际需求进行调整
            let flexWidth = 0;
            for (const char of columnContent) {
                if ((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')) {
                    // 如果是英文字符，为字符分配8个单位宽度
                    flexWidth += 8;
                } else if (char >= '\u4e00' && char <= '\u9fa5') {
                    // 如果是中文字符，为字符分配15个单位宽度
                    flexWidth += 15;
                } else {
                    // 其他种类字符，为字符分配8个单位宽度
                    flexWidth += 8;
                }
            }
            if (flexWidth < 80) {
                // 设置最小宽度
                flexWidth = 80;
            }
            if (flexWidth > 350) {
                // 设置最大宽度
                flexWidth = 350;
            }
            return flexWidth + 'px';
        };

        /**
         * 获取sql，如果有鼠标选中，则返回选中内容，否则返回输入框内所有内容
         */
        const getSql = () => {
            // 没有选中的文本，则为全部文本
            let selectSql = codemirror.value.getSelection();
            if (selectSql == '') {
                selectSql = state.sql;
            }
            return selectSql;
        };

        const saveSql = async () => {
            notEmpty(state.sql, 'sql内容不能为空');
            notNull(state.dbId, '请先选择数据库');
            await dbApi.saveSql.request({ id: state.dbId, sql: state.sql, type: 1 });
            ElMessage.success('保存成功');
        };

        /**
         * 更改数据库事件
         */
        const changeDb = (id: number) => {
            if (!id) {
                return;
            }
            clearDb();
            dbApi.tableMetadata.request({ id }).then((res) => {
                state.tableMetadata = res;
                // 赋值第一个表信息
                if (state.tableMetadata.length > 0) {
                    state.tableName = state.tableMetadata[0]['tableName'];
                    changeTable(state.tableName);
                }
            });

            dbApi.hintTables
                .request({
                    id: state.dbId,
                })
                .then((res) => {
                    state.cmOptions.hintOptions.tables = res;
                });

            dbApi.getSql.request({ id, type: 1 }).then((res) => {
                if (res) {
                    state.sql = res.sql;
                }
            });
        };

        // 清空数据库事件
        const clearDb = () => {
            state.tableName = '';
            state.tableMetadata = [];
            state.columnMetadata = [];
            state.execRes.data = [];
            state.execRes.tableColumn = [];
            state.sql = '';
            state.cmOptions.hintOptions.tables = [];
        };

        // 选择表事件
        const changeTable = async (tableName: string) => {
            if (tableName == '') {
                return;
            }
            state.columnMetadata = await dbApi.columnMetadata.request({
                id: state.dbId,
                tableName: tableName,
            });
        };

        /**
         * 自动提示功能
         */
        const showHint = () => {
            codemirror.value.showHint();
        };

        /**
         * 格式化sql
         */
        const formatSql = () => {
            let selectSql = codemirror.value.getSelection();
            // 有选中sql则只格式化选中部分，否则格式化全部
            if (selectSql != '') {
                codemirror.value.replaceSelection(sqlFormatter.format(selectSql));
            } else {
                /* 将sql内容进行格式后放入编辑器中*/
                state.sql = sqlFormatter.format(sqlFormatter.format(state.sql));
            }
        };

        const search = async () => {
            const res = await dbApi.dbs.request(state.params);
            state.dbs = res.list;
        };

        return {
            ...toRefs(state),
            cmEditor,
            changeProjectEnv,
            inputRead,
            changeTable,
            runSql,
            flexColumnWidth,
            saveSql,
            changeDb,
            clearDb,
            formatSql,
            onBeforeChange,
        };
    },
});
</script>

<style scoped lang="scss">
.codesql {
    font-size: 10pt;
    font-family: Consolas, Menlo, Monaco, Lucida Console, Liberation Mono, DejaVu Sans Mono, Bitstream Vera Sans Mono, Courier New, monospace, serif;
}
#sqlcontent {
    .CodeMirror {
        height: 300px !important;
    }
}
</style>
