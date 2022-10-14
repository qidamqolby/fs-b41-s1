const projectlist = [];
const RENDER_EVENT = "render-project";
const SAVED_EVENT = "saved-project";
const STORAGE_KEY = "project-data";

const inputProject = document.getElementById("project-input");
const inputProjectName = document.getElementById("project-name");
const inputStartDate = document.getElementById("date-start");
const inputEndDate = document.getElementById("date-end");
const inputProjectDesc = document.getElementById("project-description");
const inputUseNodeJS = document.getElementById("node-js");
const inputUseReactJS = document.getElementById("react-js");
const inputUseNextJS = document.getElementById("next-js");
const inputUseTypeScript = document.getElementById("typescript");
const inputUploadImage = document.getElementById("upload-image");

document.addEventListener("DOMContentLoaded", () => {
    inputProject.addEventListener("submit", (event) => {
        event.preventDefault();
        addProject();
        alert("Project has been added");
        inputProject.reset();
    });

    if (isStorageExist()) {
        loadDataFromStorage();
    }
});

const addProject = () => {
    const id = generateID();
    const projectName = inputProjectName.value;
    const startDate = inputStartDate.value;
    const endDate = inputEndDate.value;
    const projectDesc = inputProjectDesc.value;
    const useNodeJS = inputUseNodeJS.checked;
    const useReactJS = inputUseReactJS.checked;
    const useNextJS = inputUseNextJS.checked;
    const useTypeScript = inputUseTypeScript.checked;

    console.log(inputUploadImage.files);
    const uploadImage = URL.createObjectURL(inputUploadImage.files[0]);

    const project = {
        id,
        projectName,
        startDate,
        endDate,
        projectDesc,
        useNodeJS,
        useReactJS,
        useNextJS,
        useTypeScript,
        uploadImage,
    };

    projectlist.push(project);
    document.dispatchEvent(new Event(RENDER_EVENT));
    saveData();
};

// displaying project
document.addEventListener(RENDER_EVENT, () => {
    const listProject = document.getElementById("project-list");

    listProject.innerHTML = "";

    for (const project of projectlist) {
        const projectItem = createProjectItem(project);
        listProject.append(projectItem);
    }
});

const createProjectItem = (project) => {
    const {
        id,
        projectName,
        startDate,
        endDate,
        projectDesc,
        useNodeJS,
        useReactJS,
        useNextJS,
        useTypeScript,
        uploadImage,
    } = project;

    // count duration
    const date1 = new Date(startDate);
    const date2 = new Date(endDate);
    const diffDate = Math.abs(date2 - date1);
    const projectDuration = Math.ceil(diffDate / (1000 * 3600 * 24));

    let calculateDuration = "";
    let durationTotal = "";

    if (projectDuration > 30) {
        calculateDuration = Math.round(projectDuration / 30);
        durationTotal = `${calculateDuration} month(s)`;
    } else {
        durationTotal = `${projectDuration} day(s)`;
    }

    let createdDate = date1.getDate();
    let createdMonth = date1.getMonth();
    let createdYear = date1.getFullYear();

    const arrayMonth = [
        "January",
        "February",
        "March",
        "April",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
    ];

    for (let i = 0; i < arrayMonth.length; i++) {
        if (createdMonth - 1 === i) {
            createdMonth = arrayMonth[i];
        }
    }

    const fullDate = createdDate + " " + createdMonth + " " + createdYear;

    const projectArticle = document.createElement("article");
    projectArticle.classList.add("project-item");
    projectArticle.setAttribute("id", `${id}`);
    projectArticle.innerHTML = `
    <img src=${uploadImage} alt="">
    <div class="project-name">
        <a href="project-detail.html"><h3>${projectName}</h3></a>
        <div class="project-duration">
            <p>Duration: ${durationTotal}</p>
            <p>Created: ${fullDate}</p>
        </div>
    </div>
    <div class="project-description">
        <h3>Project Description</h3>
        <p>
            ${projectDesc}
        </p>
    </div>
    `;

    const projectTech = document.createElement("div");
    projectTech.classList.add("project-tech-info");

    if (useNodeJS) {
        const nodeJS = document.createElement("img");
        nodeJS.src = "./assets/icons/nodejs.svg";
        projectTech.append(nodeJS);
    }

    if (useReactJS) {
        const reactJS = document.createElement("img");
        reactJS.src = "./assets/icons/react-native.svg";
        projectTech.append(reactJS);
    }

    if (useNextJS) {
        const nextJS = document.createElement("img");
        nextJS.src = "./assets/icons/nextjs.svg";
        projectTech.append(nextJS);
    }

    if (useTypeScript) {
        const typescript = document.createElement("img");
        typescript.src = "./assets/icons/typescript.svg";
        projectTech.append(typescript);
    }

    const actionBtn = document.createElement("div");
    actionBtn.classList.add("action-btn");

    const editBtn = document.createElement("button");
    editBtn.classList.add("btn", "btn-primary");
    editBtn.innerText = "edit";
    editBtn.addEventListener("click", () => {
        editProject(id);
    });

    const deleteBtn = document.createElement("button");
    deleteBtn.classList.add("btn", "btn-white");
    deleteBtn.innerText = "delete";
    deleteBtn.addEventListener("click", () => {
        deleteProject(id);
    });

    actionBtn.append(editBtn, deleteBtn);

    projectArticle.append(projectTech, actionBtn);

    return projectArticle;
};

// local storage
const isStorageExist = () => {
    if (typeof Storage === undefined) {
        alert("Your browser don't support this apps");
        return false;
    }
    return true;
};

const saveData = () => {
    if (isStorageExist()) {
        const parsed = JSON.stringify(projectlist);
        localStorage.setItem(STORAGE_KEY, parsed);
        document.dispatchEvent(new Event(SAVED_EVENT));
    }
};

const loadDataFromStorage = () => {
    const localData = localStorage.getItem(STORAGE_KEY);
    let data = JSON.parse(localData);

    if (data !== null) {
        for (const project of data) {
            projectlist.push(project);
        }
    }

    document.dispatchEvent(new Event(RENDER_EVENT));
};

// generate ID
const generateID = () => {
    return +new Date();
};

// delete project
const deleteProject = (id) => {
    const projectTarget = findProjectIndex(id);

    if (projectTarget === -1) {
        return;
    }

    if (confirm("do you want to delete this project?") === true) {
        projectlist.splice(projectTarget, 1);
        document.dispatchEvent(new Event(RENDER_EVENT));
        saveData();
        window.location.reload();
    }
};

// find project index

const findProjectIndex = (projectId) => {
    for (const project of projectlist) {
        if (project.id === projectId) {
            return project;
        }
    }

    return null;
};
