import { useState } from "react";
import { BrowserRouter as Router } from "react-router-dom";
import { Navbar } from "react-bootstrap";
import {
  ConversationHeader,
  ChatContainer,
  MessageList,
  Message,
  MessageInput,
} from "@chatscope/chat-ui-kit-react";

import "./index.css";
import '@chatscope/chat-ui-kit-styles/dist/default/styles.min.css';

const App = () => {
  return (
    <Router>
      <NavBar />
      <FooterBar />
    </Router>
  );
};

const NavBar = () => {
  return (
    <Navbar expand="sm">
    </Navbar>
  );
};

type Msg = {
  message: string;
  direction: string;
}

const FooterBar = () => {
  const [messages, setMessages] = useState<Msg[]>([]);
  const onSend = (v: any) => {
    const m = {message: v, direction: "outgoing"}
    setMessages(messages => [...messages, m])
  };
  return (
    <Navbar className="footer">
      <section className="p-4 text-center w-100">
        <div className="collapse mt-3" id="collapseExample">
          <ChatContainer>
            <ConversationHeader>
              <ConversationHeader.Content userName="Emily" />
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
