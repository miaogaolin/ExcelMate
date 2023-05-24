<script>
    import { Validate, Template, OpenDefaultApp } from "../wailsjs/go/main/App";
    import {
        configData,
        outputData,
        excelData as excelStore,
        settingsData,
        templateError,
        initExcelData,
    } from "./lib/store.js";
    import { OpenExcelFile, DialogError } from "../wailsjs/go/main/App";

    import { fileBasename } from "./lib/utils";
    import Tooltip from "./component/Tooltip.svelte";

    const excelColumnName = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

    $: excelData = [];
    $: columnName = [];
    $: filename = "";
    $: basename = fileBasename(filename);
    $: {
        if (excelData.length > 0) {
            for (let i = 0; i < excelData[0].data.length; i++) {
                columnName.push(excelColumnName[i]);
            }
        }
    }

    // 接受配置的条件和模板更新
    let config = [];

    // 更新excel匹配的颜色，并且更新模板匹配的内容
    async function updateExcelColor(conf, abortSignal) {
        if (!excelData) {
            return;
        }
        outputData.update((r) => {
            r.list = [];
            return r;
        });

        for (let i = 0; i < excelData.length; i++) {
            let rowData = excelData[i].data;
            await outputResult(rowData, conf, i);
        }
    }

    async function outputResult(rowData, conf, i) {
        let matchConfigIndex = -1;
        for (let j = 0; j < conf.length; j++) {
            let res = conf[j];
            matchConfigIndex = -1;
            try {
                // 验证
                let r = await Validate(rowData, res.condition);
                if (r) {
                    matchConfigIndex = j;
                    if (res.template.trim() != "") {
                        // 模板渲染数据输出
                        let output = await Template(rowData, res.template);
                        $templateError[j] = "";
                        console.log(rowData[0]);
                        outputData.update((r) => {
                            r.list.push({
                                config_index: j,
                                text: output,
                            });
                            return r;
                        });
                    }
                    break;
                }
            } catch (e) {
                if (e.toString().indexOf("template") !== -1) {
                    if (!$templateError[j] || $templateError[j] == "") {
                        $templateError[j] = e;
                    }
                }
                break;
            }
        }
        excelData[i].config_index = matchConfigIndex;
    }

    // 选择excel文件
    function selectFile() {
        let data;
        OpenExcelFile()
            .then((res) => {
                let tmp = [];
                if (res.data && res.path != "") {
                    res.data.forEach((v) => {
                        tmp.push({
                            data: v,
                        });
                    });
                    excelData = tmp;
                    updateExcelColor(config);
                    filename = res.path;
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

    // 更新 excel 数据
    function openLocalExcel() {
        OpenDefaultApp(filename);
    }

    let tooltipText = "";
    let showTooltip = false;
    let tooltipLeft, tooltipTop;
    function colorTooltipMouseEnter(event) {
        const { clientX, clientY } = event;
        tooltipText = event.target.dataset.tooltip;
        if (tooltipText == "null") {
            return;
        } else if (tooltipText == "") {
            tooltipText = "Empty Condition";
        }

        // 计算Tooltip应该出现在哪里
        tooltipLeft = clientX + 15 + "px";
        tooltipTop = clientY + 5 + "px";

        showTooltip = true;
    }

    function colorTooltipMouseLeave() {
        showTooltip = false;
    }
    configData.subscribe((v) => {
        if (v.current == "") {
            return;
        }
        let current = getCurrentConfig(v);
        config = current && current.list ? current.list : [];
        updateExcelColor(config);
    });

    settingsData.subscribe((r) => (filename = r.excel_file));

    excelStore.subscribe((r) => {
        excelData = r;
        updateExcelColor(config);
    });
</script>

<Tooltip
    text={tooltipText}
    left={tooltipLeft}
    top={tooltipTop}
    visible={showTooltip}
/>
{#if excelData.length == 0}
    <div class="flex justify-center items-center h-full">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
            on:click={selectFile}
            class="text-indigo-500 text-lg place-content-center hover:cursor-pointer hover:text-indigo-400"
        >
            Open Excel File
        </div>
    </div>
{:else}
    <div class="py-2 fixed z-10 bg-white dark:bg-gray-800 flex items-center">
        <div class="flex">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <span
                on:click={selectFile}
                class="inline-flex items-center px-3 text-sm text-gray-900 bg-slate-100 border border-r-0 border-gray-300 rounded-l-md dark:bg-gray-600 dark:text-gray-400 dark:border-gray-600 hover:cursor-pointer flex-none"
            >
                Open File
            </span>
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div
                on:click={openLocalExcel}
                class="flex rounded-none rounded-r-lg bg-gray-50 border text-gray-900 outline-none block text-sm border-gray-300 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 flex-none hover:cursor-pointer hover:text-indigo-400"
            >
                {basename}
            </div>
        </div>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <svg
            class="fill-gray-400 dark:fill-gray-500 hover:fill-gray-700 hover:cursor-pointer dark:hover:fill-white ml-2"
            viewBox="0 0 1024 1024"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            width="28"
            height="28"
            on:click={initExcelData}
            ><path
                d="M894.481158 505.727133c0 49.589418-9.711176 97.705276-28.867468 143.007041-18.501376 43.74634-44.98454 83.031065-78.712713 116.759237-33.728172 33.728172-73.012897 60.211337-116.759237 78.712713-45.311998 19.156292-93.417623 28.877701-143.007041 28.877701s-97.695043-9.721409-142.996808-28.877701c-43.756573-18.501376-83.031065-44.98454-116.76947-78.712713-33.728172-33.728172-60.211337-73.012897-78.712713-116.759237-19.156292-45.301765-28.867468-93.417623-28.867468-143.007041 0-49.579185 9.711176-97.695043 28.867468-142.996808 18.501376-43.74634 44.98454-83.031065 78.712713-116.759237 33.738405-33.728172 73.012897-60.211337 116.76947-78.712713 45.301765-19.166525 93.40739-28.877701 142.996808-28.877701 52.925397 0 104.008842 11.010775 151.827941 32.745798 46.192042 20.977777 86.909395 50.79692 121.016191 88.608084 4.389984 4.860704 8.646937 9.854439 12.781094 14.97097l0-136.263453c0-11.307533 9.168824-20.466124 20.466124-20.466124 11.307533 0 20.466124 9.15859 20.466124 20.466124l0 183.64253c0 5.433756-2.148943 10.632151-5.986341 14.46955-3.847631 3.837398-9.046027 5.996574-14.479783 5.996574l-183.64253-0.020466c-11.307533 0-20.466124-9.168824-20.466124-20.466124 0-11.307533 9.168824-20.466124 20.466124-20.466124l132.293025 0.020466c-3.960195-4.952802-8.063653-9.782807-12.289907-14.479783-30.320563-33.605376-66.514903-60.098773-107.549481-78.753645-42.467207-19.289322-87.850837-29.072129-134.902456-29.072129-87.195921 0-169.172981 33.9533-230.816946 95.597265-61.654198 61.654198-95.597265 143.621025-95.597265 230.816946s33.943067 169.172981 95.597265 230.816946c61.643965 61.654198 143.621025 95.607498 230.816946 95.607498s169.172981-33.9533 230.816946-95.607498c61.654198-61.643965 95.597265-143.621025 95.597265-230.816946 0-11.2973 9.168824-20.466124 20.466124-20.466124C885.322567 485.261009 894.481158 494.429833 894.481158 505.727133z"
            /></svg
        >
    </div>
    <table
        class="w-full text-sm text-left text-gray-500 dark:text-gray-400 relative mt-[58px]"
    >
        <thead
            class="text-xs text-gray-700 uppercase bg-gray-100 dark:bg-gray-700 dark:text-gray-400 sticky top-[58px]"
        >
            <tr>
                <th />
                {#each columnName as c}
                    <th scope="col" class="px-6 py-3">{c}</th>
                {/each}
            </tr>
        </thead>
        <tbody>
            {#each excelData as rows}
                <tr
                    class="bg-white border-b dark:bg-gray-800 dark:border-r-gray-700 dark:border-y-gray-700 border-l-transparent"
                >
                    <td
                        class="hover:cursor-pointer"
                        data-tooltip={rows.config_index >= 0
                            ? rows.config_index +
                              1 +
                              ": " +
                              config[rows.config_index].condition
                            : "null"}
                        on:mouseenter={colorTooltipMouseEnter}
                        on:mouseleave={colorTooltipMouseLeave}
                        style="background-color:{rows.config_index >= 0
                            ? config[rows.config_index].color
                            : 'transparent'}"><div class="w-[3px]" /></td
                    >
                    {#each rows.data as r}
                        <td class="px-6 py-4">{r}</td>
                    {/each}
                </tr>
            {/each}
        </tbody>
    </table>
{/if}
