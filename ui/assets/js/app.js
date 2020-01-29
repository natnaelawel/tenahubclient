$(document).ready(function(){


var side_list = $('body .navbar-side div ul li');
side_list = side_list.slice(0,5);
var list = $('.main ul li');

list.each(function(){
        $(this).hide();
    });
    
    if($('.main ul li:nth-child(3)').hasClass('agent_active')){
        $('.main ul li:nth-child(3)').show();
    }else{
        $('.main ul li:nth-child(1)').show();
    }

side_list.css('background-color','rgb(61, 61, 61)');
side_list.eq(0).css('background-color','rgb(90, 90, 90)');

side_list.on("click", function(e){
e.preventDefault();
console.log("clicked");
side_list.css('background-color','rgb(61, 61, 61)');

var list = $('.main ul li');

list.each(function(){
    $(this).fadeOut(500);
});
var index = $(this).index() + 1;
var list = $('.main ul li:nth-child('+ index +')');

    list.fadeIn(500);
    $(this).css('background-color','rgb(80, 81, 82)');
});




///////////////////////////////////////////////////////////////////////////////////////////////
// $(document).ready(function(){
//     console.log("document is ready");
//
//     $('#selectAllCenter ').on("click", function(){
//         alert("something is occured");
//     });
// });



    // $.get("../functions/agentData.php", "", function(data,status){
    //     $("#agent_data_table").html(data);
    // });
    // $.get("../functions/healthcenterData.php", "", function(data,status){
    //     $("#healthcenter_table").html(data);
    // });
    // $.get("../functions/userData.php", "", function(data,status){
    //     $("#user_table").html(data);
    // });

    $(".selectAllCenter").click(function(){
        alert("something is occured");
    });


    checkbox1 = $('table tbody input[type="checkbox"].agent_checkbox');
    // $(document).delegate('#selectAllAgent', 'click',function(){
    //     if(this.checked){
    //         checkbox1.each(function(){
    //             this.checked = true;
    //         });
    //     } else{
    //         checkbox1.each(function(){
    //             this.checked = false;
    //         });
    //     }
    // });
    $("#selectAllAgent").click(function(){

        if(this.checked){
            checkbox1.each(function(){
                this.checked = true;
            });
        } else{
            checkbox1.each(function(){
                this.checked = false;
            });
        }
    });




        checkbox2 = $('table tbody input[type="checkbox"].center_checkbox');

$("#selectAllCenter").click(function(){


    if(this.checked){
        checkbox2.each(function(){
            this.checked = true;                        
        });
    } else{
        checkbox2.each(function(){
            this.checked = false;                        
        });
    } 
});
checkbox3 = $('table tbody input[type="checkbox"].user_checkbox');

$("#selectAllUser").click(function(){
    console.log("clickedff");
    if(this.checked){

        checkbox3.each(function(){
            this.checked = true;                        
        });
    } else{
        checkbox3.each(function(){
            this.checked = false;                        
        });
    } 
});
checkbox1.click(function(){
    if(!this.checked){
        $("#selectAll").prop("checked", false);
    }
});
checkbox2.click(function(){
    if(!this.checked){
        $("#selectAll").prop("checked", false);
    }
});
checkbox3.click(function(){
    if(!this.checked){
        $("#selectAll").prop("checked", false);
    }
});
});

/////////////////////////////////////////

// alert("center is clicked");
// $('#deleteCenterBtn').toggleClass('disabled');

///////////////////////////////////////////////////////////////////////////////////////////

    $('.navbar-side-icon').on('click', function(){
        var sidebar = $('.navbar-side');
        var main = $('.main');
        $(this).toggleClass('side-show');

        if($(this).hasClass('side-show')){
            setTimeout(function(){
                rotate(90);
                sidebar.animate({
                    left:'-30%'
                },{duration:800, queue: false});
                main.animate({
                    width:'100%',
                    left:'0%'
                },{duration:800, queue: false});
            },50);
        }else{
            setTimeout(function(){
                rotate(-180);
                 sidebar.animate({
                left:'0%',
            },{duration:800, queue: false});
            main.animate({
                width:'80%',
                left:'20%'


            },{duration:800, queue: false});
            },50);
        }
    });

    function rotate(degree) {
        var icon = $('.side-icon');
           icon.css('transform','rotate('+degree+'deg)'); 
    }



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


    $(document).on("click",'table tbody tr  td a.healthcenter-edit', function(){
        var email = $(this).closest('tr').children()[2].innerHTML;
        $("#healthcenterform .modal-body #email-healthcenter-field").val(email);

    } );
    $(document).on("click",'table tbody tr  td a.healthcenter-delete', function(){
        var email = $(this).closest('tr').children()[1].innerHTML;
        var healthcenter_id = $(this).closest('tr').children()[6].innerHTML;

        $("#deletehealthcenterform .modal-body #hidden-healthcenter-field").val(healthcenter_id);
    } );
    
    $(document).on("click",'table tbody tr  td a.user-edit', function(){
        var email = $(this).closest('tr').children()[2].innerHTML;
        var healthcenter_id = $(this).closest('tr').children()[6].innerHTML;

        $("#notifyUserForm .modal-body #email-user-field").val(email);

    } );
    $(document).on("click",'table tbody tr  td a.user-delete', function(){
        var email = $(this).closest('tr').children()[3].innerHTML;
        var user_id = $(this).closest('tr').children()[6].innerHTML;

        console.log(user_id);
        $("#deleteUserForm .modal-body #hidden-user-field").val(user_id);

    } );

    $(document).on("click",'table tbody tr  td a.agent_edit', function(){
        
        var name = $(this).closest('tr').children()[1].innerHTML;
        var firstname = name.split(" ")[0];
        var lastname = name.split(" ")[1];
        var username = $(this).closest('tr').children()[2].innerHTML;
        var email = $(this).closest('tr').children()[3].innerHTML;
        var phone = $(this).closest('tr').children()[4].innerHTML;
        var agent_id = $(this).closest('tr').children()[6].innerHTML;

        // alert(email);
        $("#editAgentForm .modal-body #firstname-field").val(firstname);
        $("#editAgentForm .modal-body #lastname-field").val(lastname);
        $("#editAgentForm .modal-body #username-field").val(username);
        $("#editAgentForm .modal-body #email-agent-field").val(email);
        $("#editAgentForm .modal-body #phone-field").val(phone);
        $("#editAgentForm .modal-body #hidden-agent-field").val(agent_id);

    } );

    $(document).on("click",'table tbody tr  td a.agent_delete', function(){
        var username = $(this).closest('tr').children()[2].innerHTML;
        var agent_id = $(this).closest('tr').children()[6].innerHTML;

        // $("#deleteAgentForm .modal-body p.deletableagent").val(username);
        // console.log(username);
        // console.log(agent_id);
        $("#deleteAgentForm .modal-body #hidden-agent-field").val(agent_id);

    } );


    $(document).on("click",'table tbody tr  td a.service_edit', function(){
        var name = $(this).closest('tr').children()[1].innerHTML;
        var description = $(this).closest('tr').children()[2].innerHTML;
        var service_id = $(this).closest('tr').children()[4].innerHTML;

        $("#editServiceForm .modal-body #service_name").val(name);
        $("#editServiceForm .modal-body #service_description").val(description);
        $("#editServiceForm .modal-body #hidden_service_id").val(service_id);


    } );
    $(document).on("click",'table tbody tr  td a.service_delete', function(){
        var service_id = $(this).closest('tr').children()[4].innerHTML;

        $("#deleteServiceForm .modal-body #hidden_service_id").val(service_id);

    } );

/////////
$('.update-alert').slideUp(2000);




