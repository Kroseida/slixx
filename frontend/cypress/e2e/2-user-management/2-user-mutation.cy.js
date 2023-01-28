describe('Manage User details', () => {

    beforeEach(() => {
        cy.loginAdmin()
    })

    it('create a new user', () => {
        cy.fixture('users').then((users) => {
            users.users.forEach((user) => {
                cy.get('#sidebar__user').click()
                cy.get('#action__create__user').click()

                cy.get('#userDetails__name').type(user.name)
                cy.get('#userDetails__firstName').type(user.firstName)
                cy.get('#userDetails__lastName').type(user.lastName)
                cy.get('#userDetails__email').type(user.email)
                cy.get('#userDetails__active').select(user.active ? 'Active' : 'Inactive')
                cy.get('#userDetails__description').type(user.description)

                cy.get('#userDetails__save__button').click()
                cy.get('.success').should('be.visible')

                user.permissions.forEach((permission) => {
                    cy.get('#permission__' + permission.replaceAll('.', '_')).click()
                });
                cy.get('#authentication__password').type(user.password)
                cy.get('#authentication__repeat_password').type(user.password)
                cy.get('#authentication__save__button').click()

                cy.get('#sidebar__dashboard').click({force: true})
            });
        })
    });

    it('update admin user', () => {
        let description = 'this is a test description for the user admin'

        cy.get('#sidebar__user').click()

        cy.get('table')
            .find('tbody')
            .find('tr')
            .first()
            .find('td')
            .first()
            .click()

        cy.get('#userDetails__description').type(description)
        cy.get('#userDetails__save__button').click()
        cy.get('.success').should('be.visible')

        cy.get('#sidebar__user').click()

        cy.get('table')
            .find('tbody')
            .find('tr')
            .first()
            .find('td')
            .first()
            .click()

        cy.get('#userDetails__description').should('have.value', 'default admin user' + description)
    });

});