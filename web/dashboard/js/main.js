// 3. connect to ws connection for new connections
// 4. render new incoming requests

document.addEventListener("DOMContentLoaded", init)

const API_REQUEST_ENDPOINT = "/__webmajor/api/request"

async function init(event) {
  const requestList = await fetchRequests()

  if (requestList === null) {
    console.log('no requests to render')
  }

  const requestListElement = document.querySelector('#request-list')
  const requestPaneElement = document.querySelector('#request-view')

  let requestStore = {}

  for (let i = 0; i < requestList.length; i++) {
    const req = requestList[i]
    requestStore[req.uuid] = req
    renderRequestListItem(requestListElement, req)
  }
}

async function fetchRequests() {
  const response = await fetch(API_REQUEST_ENDPOINT, {
    method: 'GET',
  })

  if (response.ok !== true) {
    console.log('failed to fetch requests')
    return null
  }

  return await response.json()
}

function renderRequestListItem(parent, request) {
  let item = document.createElement('div')
  item.innerHTML = `
      <div class="request-preview" data-id="${request.uuid}">
          <div class="request-preview--path">${request.method} <code>${request.requestURI}</code></div>
          <div class="request-preview--status" 
               title="${request.response.status}">${request.response.code}</div>
          <div class="request-preview--duration">${request.response.durationString}</div>
      </div>
  `

  item.addEventListener('click', function (event) {
    renderRequestPane(document.querySelector('#response-view'), request
    )
  })

  parent.appendChild(item)
}

function renderRequestPane(parent, request) {
  const renderHeaders = map => {
    let result = ''

    for (const [key, value] of Object.entries(map)) {
      result += `
        <tr>
            <td class="nowrap">${key}</td>
            <td>
                <div class="word-break bg-light p-1 rounded d-inline-block">${value}</div>
            </td>
        </tr>
      `
    }

    return result
  }

  const createdAt = new Date(request.createdAt).toLocaleTimeString()
  const reqBody = request.body === '' ? 'no body provided' : request.bodyEscaped
  const resBody = request.response.body === '' ? 'no body provided' : request.response.bodyEscaped

  parent.innerHTML = `
     <h3>${request.method} ${request.requestURI}</h3>
     <p>
       at ${createdAt}
       <span class="badge badge-primary bg-primary">${request.response.durationString}</span>
     </p>
     <table class="table table-sm">
        <thead>
        <tr><th>Header</th><th>Value</th></tr>
        </thead>
        <tbody>${renderHeaders(request.headers)}</tbody>
    </table>
    <pre class="border bg-light re-scrollable"><code>${reqBody}</code></pre>
    <hr>
    <h3>Response ${request.response.status}</h3>
    <table class="table table-sm">
        <thead>
        <tr><th>Header</th> <th>Value</th></tr>
        </thead>
        <tbody>${renderHeaders(request.response.headers)}</tbody>
    </table> 
    <pre class="border bg-light re-scrollable"><code>${resBody}</code></pre>
  `
}