import React, { Component } from 'react';
import { connect } from "react-redux";

import ChatLog from './chatlog';
import Input from './input';

function UI(props) {
  return (
    <>
      <Input {...props} />
      <ChatLog {...props} />
    </>
  );
}

export default connect(s => s)(UI);