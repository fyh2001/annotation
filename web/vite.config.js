import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from "@vitejs/plugin-vue-jsx";
import UnoCSS from "unocss/vite";
import presetIcons from "@unocss/preset-icons";
import transformerDirective from "@unocss/transformer-directives";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { NaiveUiResolver } from "unplugin-vue-components/resolvers";

// https://vitejs.dev/config/
export default defineConfig({
  base: "/annotation",
  plugins: [
    vue(),
    vueJsx(),
    UnoCSS({
      transformers: [transformerDirective()],
      presets: [presetIcons({})],
    }),
    AutoImport({
      imports: [
        "vue", { "naive-ui": ["useDialog", "useMessage", "useNotification", "useLoadingBar"] }],
      resolvers: [NaiveUiResolver()],
    }),
    Components({ resolvers: [NaiveUiResolver(),], }),
  ],
})
