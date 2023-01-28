export default {
    createClient: (token) => {
        let clientContainer = {
            nextId: 1,
            subscriptions: {},
            client: undefined,
            onConnected: () => {
            },
            onReset: () => {
            },
            onClose: () => {
            },
            listener() {
                this.client.onmessage = (event) => {
                    const data = JSON.parse(event.data);
                    if ((data.type === 'update' || data.type === 'result') && this.subscriptions[data.id].callback) {
                        this.subscriptions[data.id].callback(data);
                    } else if (data.type === 'error' && this.subscriptions[data.id] && this.subscriptions[data.id].error) {
                        this.subscriptions[data.id].error(data);
                    }
                };
            },
            subscribe(message, callback, error) {
                let id = this.nextId++;
                this.subscriptions[id] = {
                    callback,
                    error,
                };

                if (message.indexOf('mutation') !== -1) {
                    this.client.send(JSON.stringify({
                        id: id + '',
                        type: 'mutate',
                        message: {
                            query: message,
                            variables: null,
                        },
                    }));
                } else {
                    this.client.send(JSON.stringify({
                        id: id + '',
                        type: 'subscribe',
                        message: {
                            query: message,
                            variables: null,
                        },
                    }));
                }
                return id;
            },
            subscribeTrackedObject(message, callback, error) {
                let subscribeId = this.subscribe(message, (data) => {
                    if (data.message instanceof Array) {
                        this.subscriptions[subscribeId].trackedObject = data.message[0].user;
                    } else {
                        Object.entries(data.message.user).forEach(([key, value]) => {
                            this.subscriptions[subscribeId].trackedObject[key] = value;
                        });
                    }
                    callback(this.subscriptions[subscribeId].trackedObject);
                }, error);
                this.subscriptions[subscribeId].trackedObject = {};
                return subscribeId;
            },
            subscribeTrackedArray(message, callback, error) {
                let subscribeId = this.subscribe(message, (data) => {
                    if (data.message instanceof Array) {
                        data.message[0].users.forEach((user) => {
                            this.subscriptions[subscribeId].trackedObject.push(user);
                        });
                    } else {
                        Object.entries(data.message.users).forEach(([key, value]) => {
                            if (key === "$") {
                                return;
                            }
                            if (!this.subscriptions[subscribeId].trackedObject[key]) {
                                this.subscriptions[subscribeId].trackedObject.push(value[0]);
                            } else {
                                Object.entries(value).forEach(([key2, value2]) => {
                                    this.subscriptions[subscribeId].trackedObject[key][key2] = value2;
                                });
                            }
                        });
                    }
                    callback(this.subscriptions[subscribeId].trackedObject);
                }, error);
                this.subscriptions[subscribeId].trackedObject = [];

                return subscribeId;
            },
            unsubscribe(id) {
                this.client.send(JSON.stringify({
                    id: id + '',
                    type: 'unsubscribe',
                }));
                delete this.subscriptions[id];
            },
            reconnect(token) {
                if (this.client) {
                    this.client.close();
                    this.onReset();
                }
                this.client = new WebSocket(`ws://localhost:3030/graphql?authorization=${token}`);
                this.listener();
                this.client.addEventListener('open', () => {
                    this.onConnected();
                });
            }
        }
        clientContainer.reconnect(token);
        let connectedBefore = false;
        setInterval(() => {
            if (clientContainer.client.readyState === WebSocket.OPEN) {
                connectedBefore = true;
            }
            if (clientContainer.client.readyState === 3 && connectedBefore) {
                connectedBefore = false;
                clientContainer.onClose();
                clientContainer.reconnect(token);
            }
        }, 50);


        return clientContainer;
    }
}