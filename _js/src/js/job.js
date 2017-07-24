angular.module('app').controller('job', [
    '$scope',
    function ($scope) {

        $scope.removeParam = function(p){
            for (var i=0;i<$scope.params.length;i++) {
                if ($scope.params[i] == p) {
                    $scope.params.splice(i, 1);
                    break;
                }
            }
        };

        $scope.parametrized = false;
        $scope.params = [];

        $scope.addChoiceParam = function(){
            $scope.params.push({
                type: 'choice'
            });

            return false;
        };

        $scope.addStringParam = function(){
            $scope.params.push({
                type: 'string'
            });

            return false;
        };


    }
]);