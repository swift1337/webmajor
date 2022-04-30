const API_REQUEST_ENDPOINT = "/__webmajor/api/request"
const WS_REQUEST_ENDPOINT = "/__webmajor/api/ws"
const DOMAIN = document.location.host

const requestListElement = document.querySelector('#request-list')
const requestPaneElement = document.querySelector('#response-pane')
let requestStore = {}

document.addEventListener("DOMContentLoaded", init)

async function init(event) {
  const requestList = await fetchRequests()

  if (requestList === null) {
    console.log('no requests to render')
  }

  for (let i = 0; i < requestList.length; i++) {
    const req = requestList[i]
    requestStore[req.uuid] = req
    renderRequestListItem(requestListElement, req.uuid)
  }

  let requestUpdatesSocket = new WebSocket(`ws://${DOMAIN}${WS_REQUEST_ENDPOINT}`);

  requestUpdatesSocket.onmessage = function (event) {
    const req = JSON.parse(event.data)

    requestStore[req.uuid] = req

    renderRequestListItem(requestListElement, req.uuid)
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

function renderRequestListItem(parent, uuid) {
  const request = requestStore[uuid]
  let item = document.createElement('div')

  item.classList.add('request-preview')
  item.dataset.id = request.uuid

  item.innerHTML = `
    <div class="request-preview--path">${request.method} <code>${request.requestURI}</code></div>
    <div class="request-preview--status" 
         title="${request.response.status}">${request.response.code}</div>
    <div class="request-preview--duration">${request.response.durationString}</div>
  `

  item.addEventListener('click', function (e) {
    renderRequestPane(requestPaneElement, uuid)

    item.classList.toggle('active')

    const prevActiveItemElement = document.querySelector(
      `.active[data-id]:not([data-id="${request.uuid}"])`
    )

    if(prevActiveItemElement) {
      prevActiveItemElement.classList.toggle('active')
    }
  })

  parent.appendChild(item)
}

function renderRequestPane(parent, uuid) {
  const request = requestStore[uuid]
  const createdAt = new Date(request.createdAt).toLocaleTimeString()
  const reqBody = request.body === '' ? 'no body provided' : request.bodyEscaped
  const resBody = request.response.body === '' ? 'no body provided' : request.response.bodyEscaped

  parent.innerHTML = `
     <h3 class="mb-2">${request.method} ${request.requestURI}</h3>
     <p>
       <span class="badge bg-secondary">at ${createdAt}</span>
       <span class="badge bg-primary">${request.response.durationString}</span>
       <span class="badge bg-info">${request.response.status}</span>
     </p>
     <ul class="nav nav-tabs" id="responseTab" role="tablist">
       <li class="nav-item" role="presentation">
         <a class="nav-link"
            id="req-headers-tab"
            data-bs-toggle="tab"
            data-bs-target="#req-headers"
            type="button"
            role="tab"
            aria-controls="req-headers-tab"
            aria-selected="false">Request headers</a>
       </li>
       <li class="nav-item" role="presentation">
         <a class="nav-link"
            id="req-body-tab"
            data-bs-toggle="tab"
            data-bs-target="#req-body"
            type="button"
            role="tab"
            aria-controls="req-body-tab"
            aria-selected="false">Request body</a>
       </li>
       <li class="nav-item" role="presentation">
         <a class="nav-link"
            id="res-headers-tab"
            data-bs-toggle="tab"
            data-bs-target="#res-headers"
            type="button"
            role="tab"
            aria-controls="res-headers-tab"
            aria-selected="false">Response headers</a>
       </li>
       <li class="nav-item" role="presentation">
         <a class="nav-link active"
            id="res-body-tab"
            data-bs-toggle="tab"
            data-bs-target="#res-body"
            type="button"
            role="tab"
            aria-controls="res-body-tab"
            aria-selected="true">Response body</a>
       </li>
     </ul>
     <div class="tab-content" id="responseTabContent">
         <div class="tab-pane show" id="req-headers" role="tabpanel" aria-labelledby="res-headers-tab">
             <table class="table table-sm">
                <thead>
                <tr><th>Header</th><th>Value</th></tr>
                </thead>
                <tbody>${renderHeaders(request.headers)}</tbody>
            </table>
        </div>
        <div class="tab-pane" id="req-body" role="tabpanel" aria-labelledby="res-body-tab">
            <pre class="border bg-light re-scrollable"><code>${reqBody}</code></pre>
        </div>
        <div class="tab-pane" id="res-headers" role="tabpanel" aria-labelledby="res-headers-tab">
            <table class="table table-sm">
                <thead>
                <tr><th>Header</th> <th>Value</th></tr>
                </thead>
                <tbody>${renderHeaders(request.response.headers)}</tbody>
            </table> 
        </div>
        <div class="tab-pane active" id="res-body" role="tabpanel" aria-labelledby="res-body-tab">
            <pre class="border bg-light re-scrollable"><code>${resBody}</code></pre>    
        </div>
    </div>
  `
}

function renderHeaders(headers) {
  let result = ''

  for (const [key, value] of Object.entries(headers)) {
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