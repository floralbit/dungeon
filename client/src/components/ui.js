import React, { Component } from 'react';

import ChatLog from './chatlog';

class UI extends Component {
  constructor(props) {
    super(props);
    // TODO: use redux because this is awful.
    this.state = {
      messages: [],
      network: {},
    };
  }

  render() {
    return (
      <>
        <ChatLog messages={this.state.messages} network={this.state.network} />
      </>
    );
  }
}

export default UI;