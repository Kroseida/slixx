describe('Authentication via Username / Password', () => {

    beforeEach(() => {
        cy.fixture('base').then((base) => {
            cy.visit(base.url + '/#/login')
        })
    })

    it('login with username and password', () => {
        cy.get('#authentication__name').type('admin')
        cy.get('#authentication__password').type('admin')
        cy.get('#authentication__button_execute').click()

        cy.get('#localUser__name').eq(0).should('contain', 'admin')
    })

    it('login with username and invalid password', () => {
        cy.get('#authentication__name').type('admin')
        cy.get('#authentication__password').type('notAdminPassword')
        cy.get('#authentication__button_execute').click()

        cy.get('.error').should('be.visible')
    })

    it('login with invalid username and password', () => {
        cy.get('#authentication__name').type('admina')
        cy.get('#authentication__password').type('admin')
        cy.get('#authentication__button_execute').click()

        cy.get('.error').should('be.visible')
    })

});