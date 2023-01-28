import utils from './utils';

describe('Validate Permissions for Blank User', () => {

    beforeEach(() => {
        cy.fixture('users').then((users) => {
            cy.fixture('base').then((base) => {
                cy.login(base.url, users.users[0].name, users.users[0].password)
            });
        });
    })

    utils.validateUserListNotViewable()
    
})