// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
function loginAdmin() {
    cy.fixture('base').then((base) => {
        login(base.url, base.username, base.password)
    })
}

function login(url, name, password) {
    cy.visit(url + '/#/login')
    cy.get('#authentication__name').type(name)
    cy.get('#authentication__password').type(password)
    cy.get('#authentication__button_execute').click()

    cy.get('#barsTooltip').click({force: true})
}

Cypress.Commands.addAll({
    loginAdmin,
    login
})