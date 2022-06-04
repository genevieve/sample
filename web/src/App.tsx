import React from "react";
import { BrowserRouter as Router } from "react-router-dom";
import { Navbar, Container } from "react-bootstrap";
import '@chatscope/chat-ui-kit-styles/dist/default/styles.min.css';
import { MainContainer, ChatContainer, MessageList, Message, MessageInput } from "@chatscope/chat-ui-kit-react";

import "./index.css";

const App = () => {
  return (
    <Router>
      <NavBar />
      <div style={{ position:"relative", height: "500px" }}>
        <MainContainer>
          <ChatContainer>
            <MessageList>
              <Message model={{
                       message: "Hello my friend",
                       sentTime: "just now",
                       sender: "Joe"
                       }} />
              </MessageList>
            <MessageInput placeholder="Type message here" />
          </ChatContainer>
        </MainContainer>
      </div>
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
      <Container className="justify-content-center">
      </Container>
    </Navbar>
  );
};

export default App;
