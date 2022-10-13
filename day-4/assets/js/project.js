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

    const projectArticle = document.createElement("article");
    projectArticle.classList.add("project-item");
    projectArticle.setAttribute("id", `${id}`);
    projectArticle.innerHTML = `
    <img src=${uploadImage} alt="">
    <div class="project-name">
        <h3>${projectName}</h3>
        <div class="project-duration">
            <p><b>Start Date:</b> ${startDate}</p>
            <p><b>End Date:</b> ${endDate}</p>
        </div>
    </div>
    <div class="project-description">
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
