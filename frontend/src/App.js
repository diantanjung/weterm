import React, { useEffect } from "react";
import { Terminal } from "xterm";
import LocalEchoController from 'local-echo';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { AttachAddon } from 'xterm-addon-attach';
import { FitAddon } from 'xterm-addon-fit';
import { Unicode11Addon } from 'xterm-addon-unicode11';
import { SerializeAddon } from "xterm-addon-serialize";
import "./xterm.css";
import axiosInstance from "./helpers/axiosInstance";
import isAuthenticated from "./utils/isAuthenticated";

const term = new Terminal({
  convertEol: true,
  fontFamily: `Abel, monospace, MesloLGS NF`,
  fontSize: 15,
  fontWeight: 400
  // rendererType: "dom" // default is canvas
});
const localEcho = new LocalEchoController();

function App() {
  useEffect(() => {

    // Open the terminal in #terminal-container
    term.open(document.getElementById("xterm"));

    // Create a local echo controller (xterm.js >=v4)
    term.loadAddon(localEcho);

    term.focus();
    if (!isAuthenticated()) {
      startLogin();
    } else {
      startWs(localStorage.username);
    }
  }, []);

  const doLogin = (username, password) => {
    axiosInstance()
      .post("/users/login", JSON.stringify({ "username": username, password: password }))
      .then((res) => {
        if (res.data.access_token) {
          localStorage.access_token = res.data.access_token;
          localStorage.username = res.data.user.username;
          startWs(res.data.user.username);
        } else {
          term.writeln(res.data.message);
        }
      })
      .catch((err) => {
        if (err) {
          if (err.response) {
            if (err.response.data) {
              if (err.response.data.error) {
                term.writeln(err.response.data.error);
              }
            } else {
              term.writeln("Error " + err.response.status + " : " + err.response.statusText);
            }

          } else if (err.message) {
            term.writeln(err.message);
          }
        }
        startLogin();
      });
  }

  const startLogin = () => {
    // Read a single line from the user
    let loginStatus = true;
    localEcho.read("Login, username : ")
      .then(username => {
        let password = "";
        term.write("Login, password : ")
        term.onData(e => {
          if (loginStatus) {
            if (e === '\r') {
              term.write('\r\n');
              doLogin(username, password);
              password = "";
              loginStatus = false;
            } else if (e === '\u007F') {
              if (password.length > 0) {
                password = password.substr(0, password.length - 1);
              }
            } else {
              if (e >= String.fromCharCode(0x20) && e <= String.fromCharCode(0x7E)) {
                password += e;
              }
            }
          }

        });
      })
      .catch(error => alert(`Error reading: ${error}`));
  }

  const logout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("username");
    startLogin();
  }

  const startWs = (username) => {
    const url = `ws://${process.env.REACT_APP_BE_WS}/ws2/${username}`;

    const ws = new WebSocket(url);
    const attachAddon = new AttachAddon(ws);

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);

    const webLinksAddon = new WebLinksAddon();
    term.loadAddon(webLinksAddon);

    const unicode11Addon = new Unicode11Addon();
    term.loadAddon(unicode11Addon);

    const serializeAddon = new SerializeAddon();
    term.loadAddon(serializeAddon);

    ws.onclose = function (event) {
      console.log(event);
      // term.write('\r\n\nconnection has been terminated from the server-side (hit refresh to restart)\n')
     logout();
    };

    ws.onopen = function () {
      term.loadAddon(attachAddon);
      term._initialized = true;
      term.focus();
      setTimeout(function () { fitAddon.fit() });
      term.onResize(function (event) {
        var rows = event.rows;
        var cols = event.cols;
        var size = JSON.stringify({ cols: cols, rows: rows });
        var send = new TextEncoder().encode("\x01" + size);
        console.log('resizing to', size);
        ws.send(send);
      });
      term.onTitleChange(function (event) {
        console.log("evt : " + event);
        switch (event) {
          case 'open':
            window.open("https://www.google.com", "_blank");
            break;
          default:
            break;
        }
      });
      window.onresize = function () {
        fitAddon.fit();
      };
    };
  }

  return (
    <div style={{
      background: "#000000"
    }}>
      <div id="xterm" style={{ height: "100vh", width: "100%" }} />
    </div>
  );
}

export default App;
