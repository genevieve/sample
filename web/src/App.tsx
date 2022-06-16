import { useEffect, useState } from "react";
import { BrowserRouter as Router } from "react-router-dom";
import { Navbar } from "react-bootstrap";
import {
  Conversation,
  ConversationHeader,
  ConversationList,
  ChatContainer,
  MessageList,
  Message,
  MessageInput,
} from "@chatscope/chat-ui-kit-react";

import "./index.css";
import '@chatscope/chat-ui-kit-styles/dist/default/styles.min.css';

const App = () => {
  const ws = new WebSocket("ws://localhost:8080/ws?name=genevieve");

  ws.onopen = function(e) {
    console.log("[open] Connection established");
    ws.send(JSON.stringify({
      message: "genevieve",
      action: "join-room",
    }));
  };

  ws.onclose = function(event) {
    if (event.wasClean) {
      console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
      console.log('[close] Connection died');
    }
  };

  ws.onerror = function(error) {
    console.log(`[error] ${error}`);
  };

  return (
    <Router>
      <NavBar />
      <FooterBar ws={ws} />
    </Router>
  );
};

const NavBar = () => {
  return (
    <Navbar expand="sm">
    </Navbar>
  );
};

type Room = {
  id: string;
}

type Msg = {
  message: string;
  direction: string;
}

function FooterBar (props: any) {
  const [messages, setMessages] = useState<Msg[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);

  useEffect(() => {
    const unloadCallback = (event:any) => {
      event.preventDefault();
      event.returnValue = "";

      for (var i = 0; i < rooms.length; i++) {
        props.ws.send(JSON.stringify({
          message: rooms[i].id,
          action: "leave-room",
        }));
      }

      return "";
    };

    window.addEventListener("beforeunload", unloadCallback);
    return () => window.removeEventListener("beforeunload", unloadCallback);
  });


  // TODO: Extract function for receiving messages from server
  props.ws.onmessage = function(event: any) {
    let e = JSON.parse(event.data);
    console.log(`Receiving from server: ${e}`);

    if (e.action === "room-joined") {
      const r = {
        id: e.target.id,
      }
      setRooms(rooms => [...rooms, r])
    }

    if (e.action === "send-message") {
      const m = {
        message: e.message,
        direction: "incoming",
      }
      setMessages(messages => [...messages, m])
    }
  };

  const onSend = (v: any) => {
    console.log(`Sending to server: ${v}`);

    const m = {
      message: v,
      direction: "outgoing",
    }
    setMessages(messages => [...messages, m])

    props.ws.send(JSON.stringify({
      message: v,
      action: "send-message",
      target: {
        id: rooms[0].id,
      },
    }));
  };

  return (
      <Navbar className="footer">
        <section className="p-4 text-center w-100">

          <div className="collapse mt-3" id="collapseExample">
            <ConversationList>
              <Conversation name="Lilly" />
              <Conversation name="Genevieve" data-bs-toggle="collapse" href="#collapseChat1" />
            </ConversationList>

            <div className="collapse mt-3" id="collapseChat1">
              <ChatContainer>
                <ConversationHeader>
                  <ConversationHeader.Content userName="Genevieve" />
                </ConversationHeader>
                <MessageList>
                  {messages.map((item, index) => (
                  <Message
                    key={index}
                    model={{
                      message: item.message,
                      direction: item.direction,
                    }}
                  />
                  ))}
                </MessageList>
                <MessageInput
                  placeholder="Type message here"
                  onSend={(v: any) => onSend(v)}
                />
              </ChatContainer>
            </div>

          </div>

          <a className="btn btn-primary" data-bs-toggle="collapse" href="#collapseExample"
            role="button" aria-expanded="false" aria-controls="collapseExample">
            <div className="d-flex justify-content-between align-items-center">
              <span>Collapsible Chat App</span>
              <i className="fas fa-chevron-down"></i>
            </div>
          </a>

        </section>
      </Navbar>
  );
};

export default App;
