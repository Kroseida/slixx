export default (client) => ({
  /**
   * Authenticate user against the backend and return a Session.
   *
   * @param name the user name
   * @param password the user password
   * @param callback after authentication
   * @param error if authentication fails
   */
  authenticate(name, password, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "authenticate",
      args: {
        name,
        password,
      },
      fields: [
        "id",
        "token",
      ],
      isMutation: true
    });
    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(fullQuery, (data) => {
      client.graphql.unsubscribe(updatedSubscriptionId);
      callback(data);
    }, (data) => {
      error(data.message);
    });
  },
  /**
   * Subscribe to the current user. The callback is called when the user changes.
   *
   * @param subscriptionIdBefore the subscription id to unsubscribe before subscribing
   * @param callback the callback to call when the user changes
   * @param error the callback to call when an error occurs
   * @returns the subscription id
   */
  subscribeLocalUser(subscriptionIdBefore, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);

    let fullQuery = client.graphql.buildQuery({
      method: "getLocalUser",
      args: [],
      fields: [
        "id",
        "name",
        "firstName",
        "lastName",
        "email",
        "active",
        "createdAt",
        "updatedAt",
        "description",
        "permissions",
      ],
      isMutation: false
    });
    return client.graphql.subscribeTrackedObject(fullQuery, (data) => {
      try {
        callback(data);
      } catch (e) {
        error(e);
      }
    }, (data) => {
      error(data.message);
    });
  },
  subscribeUsers(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getUsers(limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows  {
          id
          name
          firstName
          lastName
          email
          active
          createdAt
          updatedAt
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  subscribeUser(subscriptionIdBefore, id, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getUser(id: "${id}") {
        id
        name
        firstName
        lastName
        email
        active
        createdAt
        updatedAt
        description
        permissions
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  createUser(args, callback, error) {
    let copiedArgs = Object.assign({}, args);
    delete copiedArgs.permission;

    let fullQuery = client.graphql.buildQuery({
      method: "createUser",
      copiedArgs,
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  },
  updateUser(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "updateUser",
      args,
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribe(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data.message[0]);
      },
      (data) => error(data.message)
    );
  },
  deleteUser(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "deleteUser",
      args: {
        id: args.id,
      },
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  },
  subscribePermissions(subscriptionIdBefore, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    let fullQuery = client.graphql.buildQuery({
      method: "getPermissions",
      args: [],
      fields: [
        "name",
        "value",
      ],
      isMutation: false
    });

    return client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        callback(data);
      },
      (data) => error(data.message)
    );
  },
  removePermission(args, callback, error) {
    let updatedSubscriptionId = client.graphql.subscribe(`
      mutation {
        data: removeUserPermission(id: "${args.id}", permissions: ["${args.permissions}"]) {
          id
        }
      }`,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data.message[0]);
      },
      (data) => error(data.message)
    );
  },
  addPermission(args, callback, error) {
    let updatedSubscriptionId = client.graphql.subscribe(`
      mutation {
        data: addUserPermission(id: "${args.id}", permissions: ["${args.permissions}"]) {
          id
        }
      }`,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data.message[0]);
      },
      (data) => error(data.message)
    );
  },
  changePassword(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "createPasswordAuthentication",
      args,
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  }
})
