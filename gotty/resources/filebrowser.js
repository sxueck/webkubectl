let App = angular.module('FileBrowserApp', []);

App.controller('FileBrowserCtrl', function ($scope, $http, $log) {
    
    $scope.currentPath = '/nonexistent';
    $scope.pathParts = [];
    $scope.files = [];
    $scope.loading = false;
    $scope.error = '';
    $scope.sessionToken = new URLSearchParams(window.location.search).get('sessionToken');

    $scope.list = function (path) {
        $scope.loading = true;
        $scope.error = '';
        $http.post('/api/filebrowser/list', { path: path, sessionToken: $scope.sessionToken })
            .then(function (response) {
                if (response.data.success) {
                    $scope.files = response.data.files;
                    $scope.currentPath = response.data.path;
                    $scope.pathParts = response.data.path.split('/').filter(p => p);
                } else {
                    $scope.error = response.data.message;
                }
                $scope.loading = false;
            }, function (response) {
                $scope.error = 'An error occurred: ' + response.statusText;
                $log.error(response);
                $scope.loading = false;
            });
    };

    $scope.refresh = function () {
        $scope.list($scope.currentPath);
    };

    $scope.navigateTo = function (path) {
        $scope.list(path);
    };

    $scope.navigateToFolder = function (folderName) {
        let newPath = ($scope.currentPath === '/' ? '' : $scope.currentPath) + '/' + folderName;
        $scope.list(newPath);
    };
    
    $scope.navigateUp = function () {
        let parts = $scope.currentPath.split('/').filter(p => p);
        parts.pop();
        let newPath = '/' + parts.join('/');
        $scope.list(newPath);
    };

    $scope.navigateToIndex = function (index) {
        let parts = $scope.currentPath.split('/').filter(p => p);
        let newPath = '/' + parts.slice(0, index + 1).join('/');
        $scope.list(newPath);
    };

    $scope.formatSize = function (size) {
        if (size < 1024) return size + ' B';
        if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB';
        if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(2) + ' MB';
        return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
    };
    
    $scope.downloadFile = function(fileName) {
        window.open(`/api/filebrowser/download?sessionToken=${$scope.sessionToken}&path=${$scope.currentPath}&file=${fileName}`, '_blank');
    };

    // Upload functionality
    $scope.selectedFile = null;
    $scope.uploading = false;
    $scope.uploadProgress = 0;
    $scope.uploadError = '';

    $scope.showUploadModal = function() {
        $('#uploadModal').modal('show');
    };

    $scope.onFileSelect = function(element) {
        $scope.$apply(function() {
            $scope.selectedFile = element.files[0];
            $scope.uploadError = '';
        });
    };

    $scope.uploadFile = function() {
        if (!$scope.selectedFile) return;

        $scope.uploading = true;
        $scope.uploadProgress = 0;
        
        let formData = new FormData();
        formData.append('file', $scope.selectedFile);
        formData.append('path', $scope.currentPath);
        formData.append('sessionToken', $scope.sessionToken);

        $http.post('/api/filebrowser/upload', formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined},
            uploadEventHandlers: {
                progress: function(e) {
                    if (e.lengthComputable) {
                        $scope.uploadProgress = Math.round((e.loaded / e.total) * 100);
                    }
                }
            }
        }).then(function(response) {
            $scope.uploading = false;
            if (response.data.success) {
                $('#uploadModal').modal('hide');
                $scope.refresh();
                $scope.selectedFile = null;
                document.getElementById('fileInput').value = '';
            } else {
                $scope.uploadError = response.data.message;
            }
        }, function(response) {
            $scope.uploading = false;
            $scope.uploadError = 'Upload failed: ' + response.statusText;
            $log.error(response);
        });
    };


    // Initial load
    $scope.list($scope.currentPath);
}); 