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
                    if ((data.type === 'update' || data.type === 'result') && this.subscriptions[data.id] && this.subscriptions[data.id].callback) {
                        this.subscriptions[data.id].callback(data);
                    } else if (data.type === 'error' && this.subscriptions[data.id] && this.subscriptions[data.id].error) {
                        this.subscriptions[data.id].error(data);
                    }
                };
            },
            buildFields(fields) {
                let query = '';
                for (let field in fields) {
                    if (typeof fields[field] === 'string') {
                        query += fields[field] + '\n';
                    }
                    if (typeof fields[field] === 'object') {
                        query += fields[field].name + '{\n' + this.buildFields(fields[field].sub) + '}\n';
                    }
                }
                return query;
            },
            buildQuery({method, args, fields, isMutation}) {
                let query = (isMutation ? 'mutation' : '') + ' {\ndata: ' + method + (args.length !== 0 ? '(' : '');
                let first = true;
                for (let key in args) {
                    if (args[key] === undefined) {
                        continue;
                    }
                    if (first) {
                        first = false;
                    } else {
                        query += ', ';
                    }
                    if (typeof args[key] === 'string') {
                        query += key + ': "' + args[key] + '"';
                    } else {
                        query += key + ': ' + args[key];
                    }
                }
                query += (args.length !== 0 ? ')' : '') + ' {\n ' + this.buildFields(fields) + '}\n}';
                return query;
            },
            subscribe(message, callback, error, id) {
                if (!id) {
                    id = this.nextId++;
                }

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
                        this.subscriptions[subscribeId].trackedObject = data.message[0].data;
                        callback(this.subscriptions[subscribeId].trackedObject);
                    } else {
                        this.unsubscribe(subscribeId)
                        this.subscribeTrackedObject(message, callback, error, subscribeId)
                    }
                }, error);
                this.subscriptions[subscribeId].trackedObject = {};
                return subscribeId;
            },
            subscribeTrackedArray(message, callback, error) {
                let subscribeId = this.subscribe(message, (data) => {
                    if (data.message instanceof Array) {
                        data.message[0].data.forEach((element) => {
                            this.subscriptions[subscribeId].trackedObject.push(element);
                        });
                    } else {
                        this.unsubscribe(subscribeId)
                        this.subscribeTrackedObject(message, callback, error, subscribeId)
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
                this.client = new WebSocket(`ws://testaaaaa@localhost:3030/graphql?authorization=${token}`);
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