import React, { Component } from 'react';
import { connect } from "react-redux";

import ChatLog from './chatlog';

function UI(props) {
  return (
    <>
      <ChatLog {...props} />
    </>
  );
}

export default connect(s => s)(UI);