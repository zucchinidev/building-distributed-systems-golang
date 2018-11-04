'use strict';

(function () {
  const $ = document.querySelector.bind(document)
  let chart
  let $options
  let $title

  function onInit () {
    $options = $('#options')
    $title = $('[data-field="title"]')
    chart = new google.visualization.PieChart($('#chart'))
    const pollId = new URL(location.href).searchParams.get('poll')
    const url = `http://localhost:8080/${pollId}?key=abc123`

    $('#delete').addEventListener('click', () => {
      if (confirm('Sure?')) {
        fetch(url, { method: 'DELETE' }).then(() => location.href = '/').catch(err => console.log(err))
      }
    })
    update(url)
  }

  function getItemListDomElements () {
    const $li = document.createElement('li')
    const $small = document.createElement('small')
    const $span = document.createElement('span')
    $small.classList.add('label', 'label-default')
    return { $small, $span, $li }
  }

  function renderResultItems ($options, key, value) {
    const { $small, $span, $li } = getItemListDomElements()
    $small.textContent = value
    $span.textContent = key
    $li.appendChild($small)
    $li.appendChild($span)
    $options.appendChild($li)
  }

  function update (url) {
    const $options = $('#options')
    const $title = $('[data-field="title"]')
    fetch(url).then(response => response.json()).then(polls => {
      const poll = polls[0]
      $options.textContent = poll.title
      $title.innerHTML = ''
      poll.options.forEach(option => {
        const value = poll.results && poll.results[option] || 0
        renderResultItems($options, option, value)
        if (poll.results) {
          const data = new google.visualization.DataTable()
          data.addColumn('string', 'Option')
          data.addColumn('number', 'Votes')
          for (const row of Object.entries(poll.results)) {
            data.addRow(row)
          }
          chart.draw(data, { is3D: true })
        }
      })
    })
    window.setTimeout(update.bind(null, url), 1000)
  }

  google.load('visualization', '1.0', { 'packages': ['corechart'] })
  google.setOnLoadCallback(function () { onInit()})
})()