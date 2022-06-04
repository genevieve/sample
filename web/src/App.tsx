import React from "react";
import { BrowserRouter as Router } from "react-router-dom";
import { Container, Navbar } from "react-bootstrap";
import '@chatscope/chat-ui-kit-styles/dist/default/styles.min.css';
import { ConversationHeader, VideoCallButton, InfoButton, ChatContainer, MessageList, Message, MessageInput } from "@chatscope/chat-ui-kit-react";

import "./index.css";

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

const FooterBar = () => {
  return (
    <Navbar className="footer">
      <section className="p-4 text-center w-100">
        <div className="collapse mt-3" id="collapseExample">
          <ChatContainer>
            <ConversationHeader>
              <ConversationHeader.Content userName="Emily" info="Active 10 mins ago" />
              <ConversationHeader.Actions>
                <VideoCallButton />
                <InfoButton />
              </ConversationHeader.Actions>
            </ConversationHeader>
            <MessageList>
              <Message model={{
                       message: "Hello my friend",
                       sentTime: "just now",
                       sender: "Joe"
                       }} />
              </MessageList>
            <MessageInput placeholder="Type message here" />
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
