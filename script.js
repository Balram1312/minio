// Function to make GET requests
async function fetchData(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching data:', error);
        return { error: error.message };
    }
}

// Function to list buckets
async function listBuckets() {
    const url = 'http://localhost:8002/buckets';
    const response = await fetchData(url);
    displayBuckets(response);
}

// Function to list objects in a bucket
async function listObjects() {
    const bucketName = document.getElementById('bucketNameInput').value;
    const url = `http://localhost:8002/objects?bucket=${encodeURIComponent(bucketName)}`;
    console.log(url)
    const response = await fetchData(url);
    displayObjects(response);
}

// Function to get presigned URL for an object
async function getPresignedUrl() {
    const fileName = document.getElementById('fileNameInput').value;
    const url = `http://localhost:8002/file?file_name=${encodeURIComponent(fileName)}`;
    const response = await fetchData(url);
    displayPresignedUrl(response);
}

// Function to handle file upload
async function uploadFile(formData) {
    try {
        const url = 'http://localhost:8002/file';
        const response = await fetch(url, {
            method: 'POST',
            body: formData
        });
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Error uploading file:', error);
        return { error: error.message };
    }
}

// Function to display buckets
function displayBuckets(response) {
    const bucketsElement = document.getElementById('buckets');
    if (response.error) {
        bucketsElement.innerHTML = `<p>Error: ${response.error}</p>`;
    } else {
        const buckets = response.buckets;
        const bucketList = buckets.map(bucket => `<li>${bucket}</li>`).join('');
        bucketsElement.innerHTML = `<ul>${bucketList}</ul>`;
    }
}

// Function to display objects
function displayObjects(response) {
    const objectsElement = document.getElementById('objects');
    if (response.error) {
        objectsElement.innerHTML = `<p>Error: ${response.error}</p>`;
    } else {
        console.log(response.data);

        const objects = response.data;
        const objectList = objects.map(object => `<li>${object}</li>`).join('');
        objectsElement.innerHTML = `<ul>${objectList}</ul>`;
    }
}

// Function to display presigned URL// Function to display presigned URL
function displayPresignedUrl(response) {
    const presignedUrlElement = document.getElementById('presignedUrl');
    if (response.error) {
        presignedUrlElement.innerHTML = `<p>Error: ${response.error}</p>`;
    } else {
        const presignedUrl = response.url;
        presignedUrlElement.innerHTML = `
            <p>Presigned URL: <a href="${presignedUrl}" target="_blank">${presignedUrl}</a></p>
            <button onclick="downloadFile('${presignedUrl}')">Download File</button>
        `;
    }
}

// Function to download file from presigned URL
function downloadFile(url) {
    const anchor = document.createElement('a');
    anchor.style.display = 'none';
    anchor.href = url;
    anchor.setAttribute('download', '');
    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
}

// Function to handle form submission (upload file)
document.getElementById('uploadForm').addEventListener('submit', async function(event) {
    event.preventDefault();
    
    const formData = new FormData(this);
    const response = await uploadFile(formData);
    displayUploadStatus(response);
});

// Function to display upload status
function displayUploadStatus(response) {
    const uploadStatusElement = document.getElementById('uploadStatus');
    if (response.error) {
        uploadStatusElement.innerHTML = `<p>Error: ${response.error}</p>`;
    } else {
        uploadStatusElement.innerHTML = `<p>Upload successful!</p>`;
    }
}


// Function to populate buckets dropdown
async function populateBuckets() {
    const url = 'http://localhost:8002/buckets';
    const response = await fetchData(url);
    const bucketSelect = document.getElementById('bucketNameSelect');
    
    if (response.error) {
        bucketSelect.innerHTML = `<option>Error: ${response.error}</option>`;
    } else {
        bucketSelect.innerHTML = response.buckets.map(bucket => `<option value="${bucket}">${bucket}</option>`).join('');
    }
}

// Function to populate objects dropdown based on selected bucket// Function to populate objects dropdown based on selected bucket
async function populateObjects() {
    const bucketName = document.getElementById('bucketNameSelect').value;
    const url = `http://localhost:8002/objects?bucket=${encodeURIComponent(bucketName)}`;
    const response = await fetchData(url);
    const objectSelect = document.getElementById('objectNameSelect');
    
    // Clear previous options
    objectSelect.innerHTML = '';

    if (response.error) {
        objectSelect.innerHTML = `<option>Error: ${response.error}</option>`;
    } else {
        objectSelect.innerHTML = response.data.map(object => `<option value="${object}">${object}</option>`).join('');
    }
}


// Function to get presigned URL for an object
async function getPresignedUrl() {
    const bucketName = document.getElementById('bucketNameSelect').value;
    const objectName = document.getElementById('objectNameSelect').value;
    
    const url = `http://localhost:8002/file?bucket=${encodeURIComponent(bucketName)}&file_name=${encodeURIComponent(objectName)}`;
    const response = await fetchData(url);

    console.log(response.url)
    displayPresignedUrl(response);
}

// Call populateBuckets() on page load to populate the buckets dropdown
populateBuckets();
