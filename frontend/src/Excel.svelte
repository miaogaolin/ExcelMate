<script>
    import { Validate, Template } from "../wailsjs/go/main/App";
    import {
        configData,
        outputData,
        excelData as excelStore,
        settingsData,
    } from "./lib/store.js";
    import { OpenExcelFile, DialogError } from "../wailsjs/go/main/App";

    import { fileBasename } from "./lib/utils";

    const excelColumnName = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

    $: excelData = [];
    $: columnName = [];
    $: basename = "";
    $: {
        let columnNum = 0;
        excelData.forEach((res) => {
            if (res.data.length > columnNum) {
                columnNum = res.data.length;
            }
        });
        for (let i = 0; i < columnNum; i++) {
            columnName.push(excelColumnName[i]);
        }
    }

    // 接受配置的条件和模板更新
    let config = [];

    // 更新excel匹配的颜色，并且更新模板匹配的内容
    async function updateExcelColor(conf) {
        if (!excelData) {
            return;
        }
        let outputTmp = [];
        for (let i = 0; i < excelData.length; i++) {
            excelData[i].color = "border-s-transparent";
            for (let j = 0; j < conf.length; j++) {
                let res = conf[j];
                try {
                    // 验证
                    let r = await Validate(excelData[i], res.condition);
                    if (r) {
                        excelData[i].color = res.color;
                        if (res.template.trim() != "") {
                            // 模板渲染数据输出
                            let output = await Template(
                                excelData[i],
                                res.template
                            );
                            outputTmp.push(output);
                        }
                        break;
                    }
                } catch (e) {
                    // console.log(e);
                }
            }
        }
        outputData.set(outputTmp);
    }

    // 选择excel文件
    function selectFile() {
        let data;
        OpenExcelFile()
            .then((res) => {
                console.log("res", res);
                let tmp = [];
                if (res.data && res.path != "") {
                    res.data.forEach((v) => {
                        tmp.push({
                            data: v,
                        });
                    });
                    excelData = tmp;
                    updateExcelColor(config);

                    // update path
                    basename = fileBasename(res.path);
                }
            })
            .catch((e) => {
                console.log("open file", e);
                DialogError(e);
            });
        return data;
    }

    function getCurrentConfig(data) {
        if (!data.files) {
            return;
        }
        for (let i = 0; i < data.files.length; i++) {
            if (data.current == data.files[i].id) {
                return data.files[i];
            }
        }
        return {};
    }

    excelStore.subscribe((r) => {
        excelData = r;
        updateExcelColor(config);
    });

    configData.subscribe((v) => {
        console.log("excel configData.subscribe", v);
        let current = getCurrentConfig(v);
        config = current && current.list ? current.list : [];
        updateExcelColor(config);
    });

    settingsData.subscribe((r) => (basename = fileBasename(r.excel_file)));
</script>

{#if excelData.length == 0}
    <div class="flex justify-center items-center h-full">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
            on:click={selectFile}
            class="text-indigo-500 hover:text-indigo-400 text-lg place-content-center hover:cursor-pointer"
        >
            Open Excel File
        </div>
    </div>
{:else}
    <div class="my-2">
        <div class="flex">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <span
                on:click={selectFile}
                class="inline-flex items-center px-3 text-sm text-gray-900 bg-slate-100 border border-r-0 border-gray-300 rounded-l-md dark:bg-gray-600 dark:text-gray-400 dark:border-gray-600 hover:cursor-pointer flex-none"
            >
                Open File
            </span>
            <input
                type="text"
                id="website-admin"
                class="flexrounded-none rounded-r-lg bg-gray-50 border text-gray-900 outline-none block min-w-0 text-sm border-gray-300 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 flex-none"
                value={basename}
                placeholder="elonmusk"
                readonly
            />
        </div>
    </div>
    <table
        class="w-full text-sm text-left text-gray-500 dark:text-gray-400 relative"
    >
        <thead
            class="text-xs text-gray-700 uppercase bg-gray-100 dark:bg-gray-700 dark:text-gray-400 sticky top-0"
        >
            <tr>
                {#each columnName as c}
                    <th scope="col" class="px-6 py-3">{c}</th>
                {/each}
            </tr>
        </thead>
        <tbody>
            {#each excelData as rows}
                <tr
                    class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 border-l-4 {rows.color}"
                >
                    {#each rows.data as r}
                        <td class="px-6 py-4">{r}</td>
                    {/each}
                </tr>
            {/each}
        </tbody>
    </table>
{/if}
