import utils from './utils';

describe('Validate Permissions for Help Desk User', () => {

    beforeEach(() => {
        cy.fixture('users').then((users) => {
            cy.fixture('base').then((base) => {
                cy.login(base.url, users.users[2].name, users.users[2].password)
            });
        });
    })

    utils.validateUserListViewable()
    utils.validateUserCreationButtonNotToBeVisible()
    utils.validateUserUpdateIsAllowed()
    utils.validateUserCreationIsBlockedByBackend()
    utils.validatePermissionUpdateIsNotAllowed()
})