const API_URL = "";


// =========================
// SAYFA GECISLERI
// =========================

function showRegister() {

    document.getElementById("loginPage").classList.add("hidden");

    document.getElementById("registerPage").classList.remove("hidden");

}


function showLogin() {

    document.getElementById("registerPage").classList.add("hidden");

    document.getElementById("loginPage").classList.remove("hidden");

}


// =========================
// KAYIT OL
// =========================

async function register() {

    const name =
        document.getElementById("registerName").value.trim();

    const email =
        document.getElementById("registerEmail").value.trim();

    const password =
        document.getElementById("registerPassword").value;

    const message =
        document.getElementById("registerMessage");


    if (!name || !email || !password) {

        message.innerText =
            "Tum alanlari doldur.";

        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/users`,
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({
                    name: name,
                    email: email,
                    password: password
                })
            }
        );


        const data =
            await response.json();


        if (!response.ok) {

            message.innerText =
                data.error || "Kayit basarisiz.";

            return;
        }


        message.innerText =
            "Kayit basarili. Giris yapabilirsin.";


        document.getElementById("registerName").value = "";

        document.getElementById("registerEmail").value = "";

        document.getElementById("registerPassword").value = "";


        setTimeout(() => {

            showLogin();

        }, 1000);


    } catch (error) {

        console.error(error);

        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// GIRIS YAP
// =========================

async function login() {

    const email =
        document.getElementById("loginEmail").value.trim();

    const password =
        document.getElementById("loginPassword").value;

    const message =
        document.getElementById("loginMessage");


    if (!email || !password) {

        message.innerText =
            "E-posta ve sifre gir.";

        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/login`,
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({
                    email: email,
                    password: password
                })
            }
        );


        const data =
            await response.json();


        if (!response.ok) {

            message.innerText =
                data.error || "Giris basarisiz.";

            return;
        }


        localStorage.setItem(
            "token",
            data.token
        );


        message.innerText =
            "Giris basarili.";


        showDashboard();


    } catch (error) {

        console.error(error);

        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// DASHBOARD
// =========================

function showDashboard() {

    document
        .getElementById("loginPage")
        .classList.add("hidden");


    document
        .getElementById("registerPage")
        .classList.add("hidden");


    document
        .getElementById("dashboardPage")
        .classList.remove("hidden");


    getFiles();

}


// =========================
// CIKIS
// =========================

function logout() {

    localStorage.removeItem("token");


    document
        .getElementById("dashboardPage")
        .classList.add("hidden");


    document
        .getElementById("loginPage")
        .classList.remove("hidden");


    document.getElementById("loginEmail").value = "";

    document.getElementById("loginPassword").value = "";

}


// =========================
// DOSYA YUKLE
// =========================

async function uploadFile() {

    const fileInput =
        document.getElementById("fileInput");

    const message =
        document.getElementById("uploadMessage");


    if (!fileInput.files.length) {

        message.innerText =
            "Lutfen bir dosya sec.";

        return;
    }


    const token =
        localStorage.getItem("token");


    if (!token) {

        message.innerText =
            "Once giris yapmalisin.";

        return;
    }


    const formData =
        new FormData();


    formData.append(
        "file",
        fileInput.files[0]
    );


    try {

        const response = await fetch(
            `${API_URL}/upload`,
            {
                method: "POST",

                headers: {
                    "Authorization":
                        `Bearer ${token}`
                },

                body: formData
            }
        );


        const data =
            await response.json();


        if (!response.ok) {

            message.innerText =
                data.error || "Dosya yuklenemedi.";

            return;
        }


        message.innerText =
            "Dosya basariyla yuklendi.";


        fileInput.value = "";


        getFiles();


    } catch (error) {

        console.error(error);

        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// DOSYALARIM
// =========================

async function getFiles() {

    const token =
        localStorage.getItem("token");


    if (!token) {
        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/files`,
            {
                method: "GET",

                headers: {
                    "Authorization":
                        `Bearer ${token}`
                }
            }
        );


        const files =
            await response.json();


        if (!response.ok) {

            console.log(files);

            return;
        }


        const fileList =
            document.getElementById("fileList");


        if (!files || files.length === 0) {

            fileList.innerHTML =
                "Henuz dosya yuklenmedi.";

            return;
        }


        fileList.innerHTML = "";


        files.forEach(file => {

            const div =
                document.createElement("div");


            div.className =
                "fileItem";


            const downloadLink =
                `${window.location.origin}/download/${file.share_token}`;


            div.innerHTML = `

                <strong>
                    ${file.file_name}
                </strong>

                <br>

                Boyut:
                ${file.file_size} byte

                <br>

                Son kullanma:
                ${new Date(file.expires_at)
                    .toLocaleString("tr-TR")}

                <div class="fileActions">

                    <button
                        onclick="window.open('${downloadLink}', '_blank')"
                    >
                        Indir
                    </button>


                    <button
                        onclick="copyLink('${downloadLink}')"
                    >
                        Linki Kopyala
                    </button>


                    <button
                        onclick="deleteFile(${file.id})"
                    >
                        Sil
                    </button>

                </div>
            `;


            fileList.appendChild(div);

        });


    } catch (error) {

        console.error(error);

    }

}


// =========================
// DOSYA SIL
// =========================

async function deleteFile(id) {

    const token =
        localStorage.getItem("token");


    if (!token) {

        alert("Once giris yapmalisin.");

        return;
    }


    const confirmDelete =
        confirm(
            "Bu dosyayi silmek istedigine emin misin?"
        );


    if (!confirmDelete) {
        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/files/${id}`,
            {
                method: "DELETE",

                headers: {
                    "Authorization":
                        `Bearer ${token}`
                }
            }
        );


        const data =
            await response.json();


        if (!response.ok) {

            alert(
                data.error ||
                "Dosya silinemedi."
            );

            return;
        }


        alert(
            data.message ||
            "Dosya basariyla silindi."
        );


        getFiles();


    } catch (error) {

        console.error(error);

        alert(
            "Sunucuya baglanilamadi."
        );

    }

}


// =========================
// LINK KOPYALA
// =========================

async function copyLink(link) {

    try {

        await navigator.clipboard.writeText(link);

        alert(
            "Paylasim linki kopyalandi."
        );

    } catch (error) {

        console.error(error);

        alert(
            "Link kopyalanamadi."
        );

    }

}


// =====================================================
// KULLANICI YONETIMI
// =====================================================


// =========================
// KULLANICILARI GETIR
// =========================

async function getUsers() {

    const token =
        localStorage.getItem("token");


    if (!token) {

        alert(
            "Once giris yapmalisin."
        );

        return;
    }


    const userList =
        document.getElementById("userList");


    userList.innerHTML =
        "Kullanicilar yukleniyor...";


    try {

        const response =
            await fetch(
                `${API_URL}/users`,
                {
                    method: "GET",

                    headers: {
                        "Authorization":
                            `Bearer ${token}`
                    }
                }
            );


        const users =
            await response.json();


        if (!response.ok) {

            userList.innerHTML =
                users.error ||
                "Kullanicilar alinamadi.";

            return;
        }


        if (!users || users.length === 0) {

            userList.innerHTML =
                "Kayitli kullanici bulunamadi.";

            return;
        }


        userList.innerHTML = "";


        users.forEach(user => {

            const div =
                document.createElement("div");


            div.className =
                "userItem";


            div.innerHTML = `

                <strong>
                    ${user.name}
                </strong>

                <br>

                E-posta:
                ${user.email}

                <br><br>

                <button
                    onclick="editUser(${user.id})"
                >
                    Guncelle
                </button>

                <button
                    onclick="deleteUser(${user.id})"
                >
                    Sil
                </button>

                <hr>

            `;


            userList.appendChild(div);

        });


    } catch (error) {

        console.error(error);


        userList.innerHTML =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// KULLANICI GETIR
// =========================

async function editUser(id) {

    const token =
        localStorage.getItem("token");


    if (!token) {

        alert(
            "Once giris yapmalisin."
        );

        return;
    }


    try {

        const response =
            await fetch(
                `${API_URL}/users/${id}`,
                {
                    method: "GET",

                    headers: {
                        "Authorization":
                            `Bearer ${token}`
                    }
                }
            );


        const user =
            await response.json();


        if (!response.ok) {

            alert(
                user.error ||
                "Kullanici alinamadi."
            );

            return;
        }


        document.getElementById(
            "updateUserId"
        ).value = user.id;


        document.getElementById(
            "updateUserName"
        ).value = user.name || "";


        document.getElementById(
            "updateUserEmail"
        ).value = user.email || "";


        document.getElementById(
            "updateUserPassword"
        ).value = "";


        document.getElementById(
            "updateUserMessage"
        ).innerText = "";


        document.getElementById(
            "updateUserBox"
        ).classList.remove("hidden");


    } catch (error) {

        console.error(error);


        alert(
            "Sunucuya baglanilamadi."
        );

    }

}


// =========================
// KULLANICI GUNCELLE
// =========================

async function updateUser() {

    const token =
        localStorage.getItem("token");


    if (!token) {

        alert(
            "Once giris yapmalisin."
        );

        return;
    }


    const id =
        document.getElementById(
            "updateUserId"
        ).value;


    const name =
        document.getElementById(
            "updateUserName"
        ).value.trim();


    const email =
        document.getElementById(
            "updateUserEmail"
        ).value.trim();


    const password =
        document.getElementById(
            "updateUserPassword"
        ).value;


    const message =
        document.getElementById(
            "updateUserMessage"
        );


    if (!id) {

        message.innerText =
            "Kullanici secilmedi.";

        return;
    }


    if (!name || !email) {

        message.innerText =
            "Ad ve e-posta zorunludur.";

        return;
    }


    const body = {
        name: name,
        email: email
    };


    if (password) {

        body.password =
            password;

    }


    try {

        const response =
            await fetch(
                `${API_URL}/users/${id}`,
                {
                    method: "PUT",

                    headers: {
                        "Content-Type":
                            "application/json",

                        "Authorization":
                            `Bearer ${token}`
                    },

                    body:
                        JSON.stringify(body)
                }
            );


        const data =
            await response.json();


        if (!response.ok) {

            message.innerText =
                data.error ||
                "Kullanici guncellenemedi.";

            return;
        }


        message.innerText =
            "Kullanici basariyla guncellendi.";


        document.getElementById(
            "updateUserPassword"
        ).value = "";


        getUsers();


        setTimeout(() => {

            cancelUpdate();

        }, 1000);


    } catch (error) {

        console.error(error);


        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// KULLANICI SIL
// =========================

async function deleteUser(id) {

    const token =
        localStorage.getItem("token");


    if (!token) {

        alert(
            "Once giris yapmalisin."
        );

        return;
    }


    const confirmDelete =
        confirm(
            "Bu kullaniciyi silmek istedigine emin misin?"
        );


    if (!confirmDelete) {
        return;
    }


    try {

        const response =
            await fetch(
                `${API_URL}/users/${id}`,
                {
                    method: "DELETE",

                    headers: {
                        "Authorization":
                            `Bearer ${token}`
                    }
                }
            );


        const data =
            await response.json();


        if (!response.ok) {

            alert(
                data.error ||
                "Kullanici silinemedi."
            );

            return;
        }


        alert(
            data.message ||
            "Kullanici basariyla silindi."
        );


        getUsers();


    } catch (error) {

        console.error(error);


        alert(
            "Sunucuya baglanilamadi."
        );

    }

}


// =========================
// GUNCELLEME IPTAL
// =========================

function cancelUpdate() {

    document.getElementById(
        "updateUserBox"
    ).classList.add("hidden");


    document.getElementById(
        "updateUserId"
    ).value = "";


    document.getElementById(
        "updateUserName"
    ).value = "";


    document.getElementById(
        "updateUserEmail"
    ).value = "";


    document.getElementById(
        "updateUserPassword"
    ).value = "";


    document.getElementById(
        "updateUserMessage"
    ).innerText = "";

}


// =========================
// SAYFA ACILINCA
// =========================

window.onload = function () {

    const token =
        localStorage.getItem("token");


    if (token) {

        showDashboard();

    }

};