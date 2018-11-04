'use strict';

(function () {
  const $ = document.querySelector.bind(document)

  function getRequestOptions () {
    const title = $('input[id="title"]').value
    const optionsValue = $('input[id="options"]').value
    const options = optionsValue.split(',').map(opt => opt.trim())
    return {
      method: 'POST',
      body: JSON.stringify({ title, options }),
      headers: { 'Content-Type': 'application/json' }
    }
  }

  function init () {
    const form = $('form#poll')
    form.addEventListener('submit', event => {
      event.preventDefault()
      const url = 'http://localhost:8080/polls/?key=abc123'
      const postOptions = getRequestOptions()
      fetch(url, postOptions).then(response => {
        location.href = `view.html?poll=${response.headers.get('Location')}`
      })
    })
  }

  window.addEventListener('load', init)
})()