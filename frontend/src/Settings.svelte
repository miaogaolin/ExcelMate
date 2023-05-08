<script>
    import { OpenConfigFile, DialogError } from "../wailsjs/go/main/App";
    import { configData, settingsData } from "./lib/store.js";

    $: path = "";

    settingsData.subscribe((r) => {
        path = r.config_file;
    });

    function selectConfigFile() {
        OpenConfigFile()
            .then((p) => {
                if (p[0] != "") {
                    path = p[0];
                    // 渲染读取的配置
                    configData.saveSet(JSON.parse(p[1]));
                }
            })
            .catch((e) => {
                DialogError(e);
            });
    }
</script>

<label
    for="website-admin"
    class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
    >Config File Path:</label
>
<div class="flex">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <span
        on:click={selectConfigFile}
        class="inline-flex items-center px-3 text-sm text-gray-900 bg-slate-100 border border-r-0 border-gray-300 rounded-l-md dark:bg-gray-600 dark:text-gray-400 dark:border-gray-600 hover:cursor-pointer"
    >
        Open File
    </span>
    <input
        type="text"
        id="website-admin"
        class="rounded-none rounded-r-lg bg-gray-50 border text-gray-900 outline-none block flex-1 min-w-0 w-full text-sm border-gray-300 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
        value={path}
        placeholder="elonmusk"
        readonly
    />
</div>
