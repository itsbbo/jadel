import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import laravel from 'laravel-vite-plugin';
import path from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
    plugins: [
        laravel({
            input: ["ui/js/app.tsx", "ui/css/app.css"],
            refresh: true,
        }),
        react({ include: /\.(mdx|js|jsx|ts|tsx)$/ }),
        tailwindcss(),
    ],
    esbuild: {
        jsx: "automatic",
    },
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./ui/js"),
            '#': path.resolve(__dirname, './ui'),
        },
    },
});