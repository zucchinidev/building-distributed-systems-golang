'use strict';

(function () {
  const $ = document.querySelector.bind(document)

  function createLink (id, title) {
    const href = `view.html?poll=polls/${id}`
    const link = document.createElement('a')
    link.setAttribute('href', href)
    link.textContent = title
    return link
  }

  function renderPolls (polls) {
    const $polls = $('#polls')
    $polls.innerHTML = ''
    for (const {id, title} of polls) {
      const link = createLink(id, title)
      const item = document.createElement('li')
      item.appendChild(link)
      $polls.appendChild(item)
    }
  }

  function showError (err) {
    console.log(err)
  }

  function update () {
    const url = 'http://localhost:8080/polls/?key=abc123'
    fetch(url).then(response => response.json()).then(renderPolls).catch(showError)
    window.setTimeout(update, 10000)
  }

  window.addEventListener('load', update)
})()