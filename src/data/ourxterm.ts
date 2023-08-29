import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
let fitAddon: FitAddon | undefined;
const term = new Terminal({
  convertEol: true,
  rows: 50,
  cols: 80
});

export function ourxterm(node: HTMLElement, data: string) {
  fitAddon = new FitAddon();
  term.loadAddon(fitAddon);
  fitAddon.fit();
  term.open(node);
  term.write(data);
}

export function handleResize() {
  if (fitAddon !== undefined) {
    fitAddon.fit();
  }
}

export function addLine(text: string){
  if(term !== undefined) {
    console.log("recieved line: "+text)
    term.write(text)
  }
}