import { cn } from '@/lib/utils';
import { FitAddon } from '@xterm/addon-fit';
import { Terminal } from '@xterm/xterm';
import { useCallback, useEffect, useRef } from 'react';

interface Props {
    onInput?: (data: string) => void;
    onOutput?: (data: string) => void;
    initialContent?: string;
}

export function useTerminal({ onInput, initialContent }: Props) {
    const containerRef = useRef<HTMLDivElement | null>(null);
    const xtermRef = useRef<Terminal | null>(null);
    const fitAddonRef = useRef<FitAddon | null>(null);

    const write = useCallback((data: string) => {
        if (xtermRef.current) {
            xtermRef.current.write(data);
        }
    }, []);

    useEffect(() => {
        if (!containerRef.current) return;

        const terminal = new Terminal({
            fontFamily: 'JetBrains Mono, monospace',
            fontSize: 14,
        });

        const fitAddon = new FitAddon();

        terminal.loadAddon(fitAddon);
        terminal.open(containerRef.current);
        fitAddon.fit();

        if (initialContent) {
            terminal.write(initialContent);
        }

        if (onInput) {
            terminal.onData((data) => {
                onInput(data);
            });
        }

        xtermRef.current = terminal;
        fitAddonRef.current = fitAddon;

        const onResize = () => fitAddon.fit();
        window.addEventListener('resize', onResize);

        return () => {
            terminal.dispose();
            window.removeEventListener('resize', onResize);
        };
    }, [containerRef, onInput, initialContent]);

    const TerminalDiv = ({ className }: { className?: string }) => <div ref={containerRef} className={cn('overflow-hidden', className)} />;

    return { TerminalDiv, write };
}
