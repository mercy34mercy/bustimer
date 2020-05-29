/// <reference path="../js/typings/jquery/jquery.d.ts" />
/// <reference path="../js/typings/jqueryui/jqueryui.d.ts" />
/// <reference path="../js/typings/custom/custom.d.ts" />
/// <reference path="../js/typings/moment/moment.d.ts" />
/// <reference path="mydiagram.ts" />
// 列挙体  --------------------------------------------------
var Mode;
(function (Mode) {
    Mode[Mode["route"] = 0] = "route";
    Mode[Mode["diagram"] = 1] = "diagram";
    Mode[Mode["busloca"] = 3] = "busloca";
    Mode[Mode["busstate"] = 4] = "busstate";
})(Mode || (Mode = {}));
var NsMode;
(function (NsMode) {
    NsMode[NsMode["normal"] = 0] = "normal"; // 通常検索
    NsMode[NsMode["quick"] = 1] = "quick"; // 今すぐ検索
    NsMode[NsMode["direct"] = 2] = "direct"; //直通検索
})(NsMode || (NsMode = {}));
var MapMode;
(function (MapMode) {
    MapMode[MapMode["pole"] = 0] = "pole"; // 駅・停留所
    MapMode[MapMode["xy"] = 1] = "xy"; // 座標
    MapMode[MapMode["polelist"] = 2] = "polelist"; // 駅・停留所（複数）
    MapMode[MapMode["direction"] = 3] = "direction"; // 経路
    MapMode[MapMode["busstate"] = 4] = "busstate"; // バス接近情報
    MapMode[MapMode["listmap"] = 5] = "listmap"; // 地図から検索
    MapMode[MapMode["busloca"] = 6] = "busloca"; // バス位置画面
})(MapMode || (MapMode = {}));
var LstBtn;
(function (LstBtn) {
    LstBtn[LstBtn["from"] = 0] = "from";
    LstBtn[LstBtn["to"] = 1] = "to";
    LstBtn[LstBtn["set"] = 2] = "set";
    LstBtn[LstBtn["decide"] = 3] = "decide";
    LstBtn[LstBtn["zero"] = 4] = "zero";
    LstBtn[LstBtn["toolarge"] = 5] = "toolarge";
})(LstBtn || (LstBtn = {}));
var PcSp;
(function (PcSp) {
    PcSp[PcSp["sp"] = 1] = "sp";
    PcSp[PcSp["pc"] = 2] = "pc";
})(PcSp || (PcSp = {}));
/// 検索メニュー
var MenuKbn;
(function (MenuKbn) {
    MenuKbn[MenuKbn["direct"] = 0] = "direct"; // 直接入力もしくは未指定
    MenuKbn[MenuKbn["map"] = 1] = "map"; // 地図
    MenuKbn[MenuKbn["word"] = 2] = "word"; // 50音
    MenuKbn[MenuKbn["eki"] = 3] = "eki"; // 鉄道・路線駅から
    MenuKbn[MenuKbn["mainstation"] = 4] = "mainstation"; // 主要駅
    MenuKbn[MenuKbn["location"] = 6] = "location"; // 観光地
    MenuKbn[MenuKbn["here"] = 7] = "here"; // 英語
    MenuKbn[MenuKbn["rosenzu"] = 9] = "rosenzu"; //路線図から
})(MenuKbn || (MenuKbn = {}));
var StationKbn = {
    Rosen: 'R',
    Bus: 'B',
    Landmark: 'L'
};
var buspackage;
(function (buspackage) {
    var define;
    (function (define) {
        define.MYDIAGRAM_STORAGE = 'MYDIAGRAMLIST';
        define.MYROUTE_STORAGE = 'MYROUTELIST';
        define.MYBUSSTATE_STORAGE = 'MYSTATELIST';
        define.HISTORY_DIAGRAM_KAISEIKBN_STORAGE = 'DIAGRAM_KAISEIKBN_STORAGE';
        define.HISTORY_INDEX_MYMENU_STORAGE = 'INDEX_MYMENY_STORAGE';
        define.HISTORY_LISTMAP_AREA_STORAGE = 'LISTMAP_AREA_STORAGE';
        define.STHISTORY_ROUTE_STORAGE = 'inputStation1';
        define.STHISTORY_DIABUSSTATE_STORAGE = 'inputStation2';
        define.EMERGENCY_HASHLIST_STORAGE = 'INFORMATIONLIST'; // 緊急メッセージのhashリストを保存するリスト
        define.INFORMARTION_SORTED_STORAGE = 'INFORMARTION_SORTED';
    })(define = buspackage.define || (buspackage.define = {}));
})(buspackage || (buspackage = {}));
/**
 * meta tag
 */
var ua = navigator.userAgent;
if ((ua.indexOf('iPhone') > 0) || ua.indexOf('iPod') > 0 || (ua.indexOf('Android') > 0 && ua.indexOf('Mobile') > 0)) {
    // スマホのとき
    $('head').prepend('<meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,minimum-scale=1.0,user-scalable=no">');
}
else {
    // PC・タブレットのとき
    $('head').prepend('<meta name="viewport" content="width=980">');
}
var jorudan = (function () {
    var nav = null;
    var info = null;
    (function ($) {
        $(onReady);
        $(window).resize(onResize);
    }(jQuery));
    function onReady() {
        var headNav = $('#headNav,.headNavClass');
        // panel-btn headnav
        $('#panel-btn').on('click', function () {
            headNav.slideToggle();
            $('#panel-btn-icon').toggleClass('close');
            return false;
        });
        // 言語選択ナビ
        var w = $(window).width();
        // グローバルナビ用
        var rwdMenu = $('#menuList'), slideSpeed = 500;
        var menuSouce = rwdMenu.html();
        // SPアコーディオンナビ用
        var menuOpen = $('.menuOpen'), menuObj = $('.menuObj');
        menuOpen.hide();
        $(window).load(function () {
            $('form').submit(Query.RemoveOtherPage);
            function menuSet() {
                if (pcsp() == PcSp.sp) {
                    if (!($('#rwdMenuWrap').length)) {
                        $('.mov-gnav').prepend('<div id="rwdMenuWrap"><div id="switchBtnArea"><a href="javascript:void(0);" id="switchBtn"><span></span><span></span><span></span></a></div></div>');
                        $('#rwdMenuWrap').append(menuSouce);
                        var menuList = $('.hideMenu');
                        // SPハンバーガー
                        $('#switchBtn').on('click', function () {
                            menuList.stop().slideToggle(slideSpeed);
                            $(this).toggleClass('btnClose');
                            if ($('.locMenu').is(':visible')) {
                                $('.locMenu').hide();
                            }
                        });
                        // SPアコーディオンナビ
                        $('.menuObj').on('click', function () {
                            $(this).toggleClass('open');
                            $(this).next('.menuOpen').stop().slideToggle();
                        });
                    }
                }
                else {
                    $('#rwdMenuWrap').remove();
                }
            }
            $(window).on('resize', function () {
                menuSet();
            });
            menuSet();
        });
        // Gナビホバー
        $('.menuObjPc').hover(function () {
            $(this).find('.menuOpen').stop(true, true).fadeIn();
        }, function () {
            $(this).find('.menuOpen').stop(true, true).fadeOut();
        });
        // SPローカルメニュー
        $('.locMenu').hide();
        if (pcsp() == PcSp.sp) {
            $('.locnavObj').on('click', function (e) {
                e.preventDefault();
                var thisObj = $($(this).attr('href'));
                if ($('.hideMenu').is(':visible')) {
                    $('.hideMenu').hide();
                }
                if (!thisObj.is(':visible')) {
                    $('.locMenu').hide();
                    $('.locnavObj').removeClass('open');
                    $(this).addClass('open');
                    thisObj.stop().slideToggle();
                    $('#switchBtn').removeClass('btnClose');
                }
                else {
                    thisObj.stop().slideUp();
                    $(this).removeClass('open');
                }
                // return false;
            });
            $('.pc').hide();
            $('.sp').show();
        }
        else {
            // PC
            $('.locnavObjPc').hover(function () {
                var thisObj = $($(this).children('a').attr('href'));
                $('.locMenu').hide();
                $('.locnavObj').removeClass('open');
                $(this).addClass('open');
                thisObj.stop().fadeIn();
            }, function () {
                var thisObj = $($(this).children('a').attr('href'));
                thisObj.stop().fadeOut();
                $(this).removeClass('open');
            });
            $('.locnavObj').on('click', function (e) {
                e.preventDefault();
            });
            $('.pc').show();
            $('.sp').hide();
        }
        //if (w > 980) {
        //    $('.setLang').hover(function () {
        //        $(this).find('ul').stop(true, true).slideDown()
        //    }, function () {
        //        $(this).find('ul').stop(true, true).slideUp()
        //    });
        //}
        //// オプションのメニュー表示
        //if (w < 980) {
        //    $('.optionTitle').on('click', function () {
        //        $(this).next('ul').slideToggle();
        //        return false;
        //    });
        //}
        // パン屑削除
        if (pcsp() == PcSp.sp) {
            $('.pankuzu').hide();
            $('.jbacklnk').show();
            $('.spmelogo').show();
        }
        else {
            $('.pankuzu').show();
            $('.jbacklnk').hide();
            $('.spmelogo').hide();
        }
        // 交通選択のメニュー表示
        //$('.clickMenu,.historyMenu').hide();
        $('.historyMenu').hide();
        var over_flg = true;
        $('.trafficBtn').on('click', function () {
            if ($(this).hasClass('selected')) {
                // メニュー非表示
                $(this).removeClass('selected').next('ul').slideUp('fast');
            }
            else {
                // 表示しているメニューを閉じる
                $('.trafficBtn').removeClass('selected');
                //$('.clickMenu,.historyMenu').hide();
                $('.historyMenu').hide();
                // メニュー表示
                $(this).addClass('selected').next('ul').slideDown('fast');
            }
            return false;
        });
        // マウスカーソルがメニュー上/メニュー外
        //$('.trafficBtn,.clickMenu,.historyMenu').hover(function () {
        $('.trafficBtn,.historyMenu').hover(function () {
            over_flg = true;
        }, function () {
            over_flg = false;
        });
        // メニュー領域外をクリックしたらメニューを閉じる
        $('body').on('click', function () {
            if (over_flg == false) {
                $('.trafficBtn').removeClass('selected');
                //$('.clickMenu,.historyMenu').slideUp('fast');
                $('.historyMenu').slideUp('fast');
            }
        });
        // セレクトボックス装飾
        window.addEventListener('load', function () {
            var selectArea = document.querySelectorAll('.selectArea');
            selectView(selectArea);
        }, false);
        function selectView(elm) {
            var selectObj = new Array();
            var selectTxt = new Array();
            for (var i = 0, len = elm.length; i < len; i++) {
                selectObj.push(elm[i].querySelector('.selectObj'));
                selectTxt.push(elm[i].querySelector('.selectTxt'));
                //初期表示
                selectTxt[i].innerHTML = selectObj[i].options[selectObj[i].selectedIndex].text;
                //プルダウン変更時
                selectObj[i].onchange = (function (i) {
                    return function () {
                        selectTxt[i].innerHTML = this.options[this.selectedIndex].text;
                    };
                })(i);
            }
        }
        //// セレクトボックス装飾
        //window.addEventListener(
        //    'load',
        //    function () {
        //        var selectArea = document.querySelectorAll('.selectArea');
        //        selectView(selectArea);
        //        // 運行情報イベント設定
        //        ////setInterval(GetInfomation, 60 * 1000);
        //        ////GetInfomation();
        //    },
        //    false
        //);
        //function selectView(elm) {
        //    var selectObj = new Array();
        //    var selectTxt = new Array();
        //    for (var i = 0, len = elm.length; i < len; i++) {
        //        selectObj.push(elm[i].querySelector('.selectObj'));
        //        selectTxt.push(elm[i].querySelector('.selectTxt'));
        //        //初期表示
        //        var option = selectObj[i].options[selectObj[i].selectedIndex];
        //        if (typeof option != 'undefined') {
        //            selectTxt[i].innerHTML = option.text;
        //        }
        //        //プルダウン変更時
        //        selectObj[i].onchange = (function (i) {
        //            return function () {
        //                selectTxt[i].innerHTML = this.options[this.selectedIndex].text;
        //            }
        //        })(i);
        //    }
        //}
        //$(function () {
        //	// マップのバス停など表示
        //	$(".customCheck input[type='checkbox']").on('change',function(){
        //		var thisClass = $(this).next('label').attr('for');
        //		if($(this).is(":checked")){
        //			$('.' + thisClass).show();
        //		}else{
        //			$('.' + thisClass).hide();
        //		}
        //	});
        //	$('.popUp').on('click',function(){
        //		if($('.pointerArea').is(':visible')){
        //			$('.pointerArea').hide();
        //		}else{
        //			$('.pointerArea').show();
        //		}
        //	});
        //});
        // pagatop
        var btnTop = $('#pagetop');
        btnTop.hide();
        $(window).on('scroll', function () {
            // お知らせガイド実施中は非表示にする
            if (!(nav && nav.IsExecute) && $(this).scrollTop() > 200) {
                btnTop.fadeIn();
            }
            else {
                btnTop.fadeOut();
            }
        });
        btnTop.on('click', function () {
            $('html,body').animate({ scrollTop: 0 }, 500);
            return false;
        });
        if ($('#usemore').val() == "1") {
            readMore();
            // 詳細を見るエリアの文字数丸め処理
            function readMore() {
                if (typeof $('.readMore').jTruncSubstr != 'undefined') {
                    $('.readMore').jTruncSubstr({
                        length: 80,
                        minTrail: 0,
                        moreText: "詳細を見る▼",
                        lessText: "閉じる▲",
                        ellipsisText: " ...",
                        moreAni: 0,
                        lessAni: 0
                    });
                }
            }
        }
        else {
            $.each($('.readMore'), function (idx, el) {
                var readmore = $(el).html();
                $(el).html(readmore.replace(new RegExp('&lt;', "g"), '<').replace(new RegExp('&gt;', "g"), '>'));
            });
        }
        //del 2016/02/05 m.ohkoshi del spot.jsへ移動
        //// ドロップダウンメニュー
        //var acMenulist = $('.spotList li'),
        //closeMenu  = $('.closeMenu');
        //closeMenu.hide();
        //acMenulist.on('click',function(e) {
        //	$(this).children('a').toggleClass('active');
        //	$(this).children('ul.closeMenu').stop().slideToggle();
        //	$(this).siblings('li').children('ul.closeMenu').stop().slideUp();
        //	$(this).siblings('li').children('a').removeClass('active');
        //	e.stopPropagation();
        //});
        //placeholder対応
        //$(function () {
        var supportsInputAttribute = function (attr) {
            var input = document.createElement('input');
            return attr in input;
        };
        if (!supportsInputAttribute('placeholder')) {
            $('[placeholder]').each(function () {
                var input = $(this), placeholderText = input.attr('placeholder'), placeholderColor = 'GrayText', defaultColor = input.css('color');
                input.
                    focus(function () {
                    if (input.val() === placeholderText) {
                        input.val('').css('color', defaultColor);
                    }
                }).
                    blur(function () {
                    if (input.val() === '') {
                        input.val(placeholderText).css('color', placeholderColor);
                    }
                    else if (input.val() === placeholderText) {
                        input.css('color', placeholderColor);
                    }
                }).
                    blur().
                    parents('form').
                    submit(function () {
                    if (input.val() === placeholderText) {
                        input.val('');
                    }
                });
            });
        }
        var elm = document.querySelectorAll('.selectArea');
        var selectObj = new Array();
        var selectTxt = new Array();
        for (var i = 0, len = elm.length; i < len; i++) {
            selectObj.push(elm[i].querySelector('.selectObj'));
            selectTxt.push(elm[i].querySelector('.selectTxt'));
            //初期表示
            var option = selectObj[i].options[selectObj[i].selectedIndex];
            if (typeof option != 'undefined') {
                selectTxt[i].innerHTML = option.text;
            }
            //プルダウン変更時
            selectObj[i].onchange = (function (i) {
                return function () {
                    selectTxt[i].innerHTML = this.options[this.selectedIndex].text;
                };
            })(i);
        }
        // 運行情報イベント設定
        info = new Information();
        //info.AddTabEvent($('.bodyInfo'));
        info.SortInformation(pcsp() == PcSp.pc ? '.infoAreaPcTab,.infoAreaPc' : '.infoAreaSp');
        // ソートダイアログ
        if ($('.openSortDialog').length > 0) {
            $('.openSortDialog').on('click', function () {
                info.OpenSortDialogEvent(pcsp() == PcSp.pc ? '.infoAreaPcTab,.infoAreaPc' : '.infoAreaSp');
            });
        }
        // 使い方ガイド
        AddNavEvent($('.btnNav'), $('#navparams').val());
        // My時刻表
        if ($('.mydiagram').length > 0) {
            new MyDiagram().getDiagram('.mydiagram');
        }
        // 緊急情報ポップアップイベント
        if ($('.hasEmergency').length > 0) {
            info.EmergetncyDialogEvent($('.hasEmergency'), $('.evergencyBtnArea'));
        }
        //});
        // onReady
    }
    function onResize() {
        // 言語選択ナビ
        if (pcsp() == PcSp.pc) {
            $('.setLang').hover(function () {
                $(this).find('ul').stop(true, true).slideDown();
            }, function () {
                $(this).find('ul').stop(true, true).slideUp();
            });
            // パン屑削除
            if (pcsp() == PcSp.sp) {
            }
            else {
                $('.pankuzu').show();
                $('.jbacklnk').hide();
                $('.spmelogo').hide();
            }
            $('.pc').show();
            $('.sp').hide();
        }
        else {
            $('#selectLang').css('display', 'block');
            $('.pankuzu').hide();
            $('.jbacklnk').show();
            $('.spmelogo').show();
            $('.pc').hide();
            $('.sp').show();
        }
        // onResize
    }
    /**
     * 使い方ガイド実装
     * @param $tabdiv
     */
    function AddNavEvent($navButton, navParams) {
        if (!navParams) {
            return;
        }
        if (nav == null) {
            nav = new Nav(navParams);
        }
        nav.AddBaloon();
        $navButton.on('click', nav.Show);
    }
}());
/**
 * お知らせ
 */
var Information = (function () {
    function Information() {
    }
    ///**
    // * タブを切り替えた時に表示を切り替えるイベント（お知らせ等）
    // * @param $tabdiv
    // */
    //public AddTabEvent($tabdiv: JQuery) {
    //    $tabdiv.children('nav').children('a').on('click', function () {
    //        // 選択タブに「active」追加
    //        $(this).siblings().removeClass('active');
    //        $(this).addClass('active');
    //        // 対応するタブにactive追加など
    //        var infoid = $(this).attr('data-information-id');
    //        var $infoboxes = $tabdiv.find('.infoBox');
    //        $infoboxes.removeClass('active none');
    //        if (pcsp() == PcSp.pc) {
    //            $infoboxes.filter('#' + infoid).addClass('active');
    //            $infoboxes.not('#' + infoid).addClass('none');
    //        }
    //    });
    //}
    //public ResizeEvent = (e: JQueryEventObject) => {
    //}
    Information.prototype.SetInformation = function (_more, _close) {
        // お知らせタブ切り替え
        var tab = $('.info .tab');
        $('.info .tab01').addClass('current');
        tab.on('click', function () {
            var index = tab.index(this);
            $('.infoConts').css({ 'display': 'none' });
            $('.infoConts').eq(index).css({ 'display': 'block' });
            tab.removeClass('current');
            $(this).addClass('current');
            return false;
        });
        // 詳細を見るエリアの文字数丸め処理
        function readMore() {
            $('.readMore').jTruncSubstr({
                length: 34,
                minTrail: 0,
                moreText: _more + "▼",
                lessText: _close + "▲",
                ellipsisText: " ...",
                moreAni: 0,
                lessAni: 0
            });
        }
        //PCSP共通　お知らせ内のテキスト詳細を見る
        if ($('#usemore').val() == "1") {
            readMore();
        }
        else {
            $.each($('.readMore'), function (idx, el) {
                var readmore = $(el).html();
                $(el).html(readmore.replace(new RegExp('&lt;', "g"), '<').replace(new RegExp('&gt;', "g"), '>'));
            });
        }
        // 詳細を見る
        $('.moreArea').hide();
        $('.detailMore').on('click', function () {
            $(this).parent('.mytoolTx').next('.mytoolTx').show();
            $(this).hide();
            return false;
        });
        // 閉じる
        $('.detailClose').on('click', function () {
            $(this).parent('.mytoolTx').prev('.mytoolTx').find('.detailMore').show();
            $(this).parent('.mytoolTx').hide();
            return false;
        });
        //SP時　お知らせ　横矢印アイコンのリスト開閉
        var infoTtl = $('.infoTtl');
        infoTtl.on('click', function () {
            $(this).next().stop().slideToggle();
            $(this).toggleClass('open');
        });
    };
    /**
     * 緊急情報(header)のダイアログ表示イベント
     */
    Information.prototype.EmergetncyDialogEvent = function ($emergencyArea, $btn) {
        var _this = this;
        //hiddenからlist作成
        var emergencyList = this.getEmergencyDataList($emergencyArea);
        $btn.on('click', function () { _this.openEmergencyDialog(emergencyList); });
        //setInterval(function () {
        //    $emergencyArea.toggleClass('notview');
        //}, 1000);
        // スマホの場合、確認済の緊急リストhashの確認
        if (pcsp() == PcSp.sp) {
            if (!localStorage) {
                return;
            }
            // localStorage取得
            var emergencyHash = localStorage.getItem(buspackage.define.EMERGENCY_HASHLIST_STORAGE);
            var emergencyHashList;
            if (emergencyHash == null) {
                emergencyHashList = [];
            }
            else {
                emergencyHashList = emergencyHash.split(',');
            }
            // emergencyListのhashをループし、emergencyHashListにあるかどうかの判定
            var isOpenDialog = false;
            for (var i = 0; i < emergencyList.length; i++) {
                if ($.inArray(emergencyList[i]['hash'], emergencyHashList) == -1) {
                    isOpenDialog = true;
                    emergencyHashList.push(emergencyList[i]['hash']);
                }
            }
            if (isOpenDialog) {
                localStorage.setItem(buspackage.define.EMERGENCY_HASHLIST_STORAGE, emergencyHashList.join(','));
                this.openEmergencyDialog(emergencyList);
            }
        }
    };
    Information.prototype.getEmergencyDataList = function ($emergencyArea) {
        //hiddenからlist作成
        var emergencyList = [];
        $emergencyArea.eq(0).find('.emergencyData').each(function (index, elem) {
            emergencyList.push({
                'announce_datetime': $(elem).find('.announce_datetime').val(),
                'title': $(elem).find('.title').val(),
                'message': $(elem).find('.message').val(),
                'hash': $(elem).find('.hash').val()
            });
        });
        return emergencyList;
    };
    Information.prototype.openEmergencyDialog = function (emergencyList) {
        if ($('#emergencyDialog').length > 0) {
            $('#emergencyDialog').dialog();
            return;
        }
        var $dialog = $('<div/>', {
            id: 'emergencyDialog', 'class': 'emergencyArea'
        }).appendTo('body');
        var $inner = $('<div/>', {
            id: 'emergencyDialogInner', 'class': 'emergencyAreaInner'
        }).appendTo($dialog);
        // 高さを設定(SP)
        $inner.css('max-height', $(window).height() * 0.6 + 'px');
        // 取得した緊急情報でループ
        for (var i = 0; i < emergencyList.length; i++) {
            var $dl = $('<dl/>', {
                'class': 'infoTxt'
            }).appendTo($inner);
            // 日付
            $('<dt/>', {
                'class': 'infoDate',
                text: emergencyList[i]['announce_datetime']
            }).appendTo($dl);
            var $dd = $('<dd/>', {
                'class': 'readMore'
            }).appendTo($dl);
            // タイトル
            $('<span/>', {
                'class': 'infoLine',
                text: emergencyList[i]['title']
            }).appendTo($dd);
            // メッセージ
            $('<div/>', {
                html: emergencyList[i]['message']
            }).appendTo($dd);
            // hash
            $('<input/>', {
                value: emergencyList[i]['hash'],
                type: 'hidden',
                'class': 'hash'
            }).appendTo($dd);
        }
        var buttons = [
            {
                text: $('.emergencyClose').val(),
                click: function () {
                    $('#emergencyDialog').dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        $dialog.dialog({
            modal: true,
            buttons: buttons,
            dialogClass: 'msgdialog',
            title: $('.emergencyTitle').val(),
            close: function (event) {
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            }
        });
    };
    /**
     * お知らせを並べ替え
     * @param _selector
     */
    Information.prototype.SortInformation = function (_tabSelector) {
        var $tabdiv = $(_tabSelector);
        // localStorage未対応であれば終了
        if (!localStorage) {
            return;
        }
        var informations = localStorage.getItem(buspackage.define.INFORMARTION_SORTED_STORAGE);
        // 未登録であれば終了
        if (!informations) {
            return;
        }
        var informationlist = informations.split(',');
        // 会社コードリストをループし、末尾に移動する
        for (var i = 0; i < informationlist.length; i++) {
            // タブのヘッダー・中身両方にループ
            $tabdiv.each(function (index0, elem) {
                var $ul = null;
                // PCトップページ等、標準の場合
                if ($(elem).find('ul').length > 0) {
                    $ul = $(elem).find('ul');
                }
                else if ($(elem).prop("tagName").toLowerCase() == 'ul') {
                    $ul = $(elem);
                }
                if ($ul != null) {
                    $ul.each(function (index1, elem1) {
                        $(elem1).find('.company_code').each(function (index2, elem2) {
                            // 会社コードが合致していれば、その親要素に代入
                            if (informationlist[i] == $(elem2).val()) {
                                $(elem2).parents('li').appendTo($(elem1));
                            }
                        });
                    });
                }
                else {
                    $(elem).find('dl').each(function (index1, elem1) {
                        $(elem1).find('.infoTil').find('.company_code').each(function (index2, elem2) {
                            // 会社コードが合致していれば、その親要素に代入
                            if (informationlist[i] == $(elem2).val()) {
                                $(elem2).closest('.infoArea').appendTo($(elem));
                            }
                        });
                    });
                }
            });
        }
    };
    /**
     * 並べ替えのダイアログを開く
     * @param _selector
     */
    Information.prototype.OpenSortDialogEvent = function (_tabSelector) {
        var _this = this;
        var $tabdiv = $(_tabSelector);
        if ($tabdiv.length == 0) {
            return;
        }
        if ($('#infoSortDialog').length > 0) {
            $('#infoSortDialog').dialog();
            return;
        }
        var $dialog = $('<div/>', {
            id: 'infoSortDialog'
        }).appendTo('body');
        // メッセージ
        $('<p/>', {
            html: Query.GetParameter('message', $('#infosortlabels').val())
        }).appendTo($dialog);
        // 会社名・会社コードリスト取得
        var companylist = [];
        var $ul2 = null;
        // PCトップページ等、標準の場合
        if ($tabdiv.find('ul').length > 0) {
            $ul2 = $tabdiv.find('ul');
        }
        else if ($tabdiv.prop("tagName").toLowerCase() == 'ul') {
            $ul2 = $tabdiv;
        }
        if ($ul2 != null) {
            $ul2.eq(0).children('li').each(function (index, elem) {
                // 緊急情報はスキップ
                if (parseBool($(elem).find('.isEmergency').val())) {
                    return true;
                }
                companylist.push({ company_code: $(elem).find('.company_code').val(), company_name: $(elem).find('.company_name').val().trim() });
            });
        }
        else {
            $tabdiv.children('dl').each(function (index, elem) {
                // 緊急情報はスキップ
                if (parseBool($(elem).find('.isEmergency').val())) {
                    return true;
                }
                companylist.push({
                    company_code: $(elem).find('.company_code').val(), company_name: $(elem).find('.company_name').val().trim()
                });
            });
        }
        var $ul = $('<ul/>', {
            'class': 'infoSortArea'
        }).appendTo($dialog);
        // 取得した会社名・会社コードリストでループ
        for (var i = 0; i < companylist.length; i++) {
            var $span = $('<span/>', {
                text: companylist[i]['company_name']
            });
            // リストに追加
            $ul.append($('<li/>', {
                'data-companycode': companylist[i]['company_code'],
                'class': 'company_code_' + companylist[i]['company_code'],
                html: $span
            }));
        }
        var buttons = [
            {
                text: Query.GetParameter('savebutton', $('#infosortlabels').val()),
                click: function () {
                    // 保存
                    var informationlist = $.map($('.infoSortArea').children('li'), function (v) { return v.getAttribute('data-companycode'); });
                    localStorage.setItem(buspackage.define.INFORMARTION_SORTED_STORAGE, informationlist.join(','));
                    // 再度並べ替え
                    _this.SortInformation(_tabSelector);
                    $('#infoSortDialog').dialog('close');
                },
                class: 'mouseon btn'
            },
            {
                text: Query.GetParameter('cancelbutton', $('#infosortlabels').val()),
                click: function () {
                    $('#infoSortDialog').dialog('close');
                },
                class: 'mouseon btn2'
            }
        ];
        // ソート実装
        $ul.sortable();
        $ul.disableSelection();
        $dialog.dialog({
            modal: true,
            buttons: buttons,
            dialogClass: 'msgdialog',
            title: Query.GetParameter('title', $('#infosortlabels').val()),
            close: function (event) {
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            }
        });
    };
    return Information;
}());
/**
 * ナビ（使い方ガイド）
 */
var Nav = (function () {
    function Nav(navParams) {
        var _this = this;
        this.navlist = [];
        this.IsExecute = false;
        this.showNo = 0;
        /**
         * ナビゲーション表示
         */
        this.Show = function (ev) {
            // 実行中なら終了
            if (_this.IsExecute) {
                return;
            }
            if (_this.navlist.length == 0) {
                return;
            }
            _this.IsExecute = true;
            _this.showNo = 0;
            _this.showDisp(0);
            // 画面クリック時イベント（ナビ終了）
            ev.stopPropagation();
            $(window).on('click', _this.Stop);
        };
        this.Stop = function (ev) {
            if (_this.IsExecute) {
                $('.navindex').remove();
            }
            _this.IsExecute = false;
            $(window).off('click', _this.Stop);
        };
        this.navlist = [];
        try {
            var navParamJson = JSON.parse(navParams);
            for (var key in navParamJson) {
                this.navlist.push({
                    selector: navParamJson[key].selector,
                    title: navParamJson[key].title,
                    message: navParamJson[key].message,
                    isnav: parseBool(navParamJson[key].isnav),
                    isbaloon: parseBool(navParamJson[key].isbaloon)
                });
            }
        }
        catch (e) {
        }
    }
    /**
     * バルーン（オンマウス時イベント）追加
     */
    Nav.prototype.AddBaloon = function () {
        if (this.navlist.length == 0 || pcsp() == PcSp.sp) {
            return;
        }
        for (var i = 0; i < this.navlist.length; i++) {
            var nav = this.navlist[i];
            if (!nav.isbaloon) {
                continue;
            }
            if ($(nav.selector).length == 0) {
                continue;
            }
            // 画面からはみ出そうな場合
            var position;
            if (pcsp() == PcSp.pc) {
                position = $(nav.selector).offset().left + $(nav.selector).outerWidth() + 300 > $(window).width() ? "bottom left" : "bottom right";
            }
            else {
                var scrollPosition = $(nav.selector).offset().top - $(window).scrollTop(); // 画面上の高さ取得
                if ($(window).height() / 2 < scrollPosition) {
                    position = "top";
                }
                else {
                    position = "bottom";
                }
            }
            $(nav.selector).balloon({
                html: true,
                contents: this.getMessageHtml(nav),
                delay: 1000,
                position: position,
                classname: 'navmessage balloon',
                css: null,
                showDuration: 200,
                hideDuration: 200,
                showAnimation: function (d) { this.fadeIn(d); }
            });
        }
    };
    /**
     * 画面にメッセージをfadein・fadeout
     * @param index
     */
    Nav.prototype.showDisp = function (index) {
        var _this = this;
        var nav = this.navlist[index];
        var $navTarget = $(nav.selector);
        var fadecount = 0;
        var currentOffset = null;
        var currentWidth = 0;
        // 表示対象の要素が1件も無ければ終了
        // もしくは画面から表示されていなければ終了
        if ($navTarget.length == 0 || !nav.isnav || !$navTarget.is(':visible')) {
            // 次があればshowDisp再帰実行
            if (this.navlist.length > (index + 1)) {
                this.showDisp(index + 1);
            }
            else {
                this.IsExecute = false;
            }
            return;
        }
        this.showNo++;
        // 対象となるクラスだけ、■で囲む
        for (var i = 0; i < $navTarget.length; i++) {
            // 表示されていなければcontinue
            if (!$navTarget.eq(i).is(':visible')) {
                continue;
            }
            var elem = $navTarget[i];
            var offset = $(elem).offset();
            //// ■用のDiv作成
            var $div = $('<div/>', {
                'class': 'navdiv navindex navindex' + index
            });
            $div.offset({ top: (offset.top - 8), left: (offset.left - 8) }); // 余白5 + 線太さ3
            $div.width($(elem).outerWidth() + 10);
            $div.height($(elem).outerHeight() + 10);
            $div.appendTo($('body'));
            ///// 番号作成
            if (pcsp() == PcSp.pc) {
                var $divno = $('<div/>', {
                    'class': 'navno navindex navindex' + index
                });
                $divno.appendTo($('body'));
                $divno.offset({ top: (offset.top - $divno.outerHeight() / 2), left: (offset.left + $div.width() - $divno.outerWidth() / 2) });
                $divno.append($('<p/>', { 'text': 'STEP' })).append($('<p/>', { 'text': this.showNo, 'class': 'navnono' }));
            }
            if (i == 0) {
                currentOffset = offset;
                currentWidth = $div.width();
            }
        }
        // メッセージ作成
        var $divmessage = $('<div/>', {
            'class': 'navmessage navindex navindex' + index
        });
        $divmessage.html(this.getMessageHtml(nav));
        $divmessage.appendTo($('body'));
        ///// 表示位置の調整
        // PC版の場合の画面位置
        if (pcsp() == PcSp.pc) {
            var currentleft = currentOffset.left + currentWidth + 60;
            // 画面からはみ出そうな場合
            if (currentleft + 300 > $(window).width()) {
                currentleft = currentOffset.left - 300 - 60;
            }
            //rev 2016-12-19 m.ohkoshi ---s
            var currenttop = currentOffset.top - $divmessage.eq(0).height() - 30;
            if (currentleft < 0) {
                currentleft = 0;
            }
            //rev 2016-12-19 m.ohkoshi ---e
            $divmessage.offset({
                top: currenttop, left: currentleft
            });
        }
        else {
            // 幅設定
            $divmessage.outerWidth($(window).width());
            //scrollTo
            $("html,body").animate({ scrollTop: currentOffset.top - 100 });
            var scrollPosition = currentOffset.top - $(window).scrollTop(); // 要素のtopから、画面のtop引く
            // 要素が画面の下半分にある時には、上側にウィンドウメッセージ表示
            if ($(window).height() / 2 < scrollPosition) {
                $divmessage.offset({
                    top: currentOffset.top - $divmessage.eq(0).height() - 30,
                    left: 0
                });
            }
            else {
                $divmessage.offset({
                    top: currentOffset.top + $navTarget.eq(0).height() + 30,
                    left: 0
                });
            }
        }
        // フェードイン対象のカウンタ
        fadecount = $('.navindex' + index).length;
        $('.navindex' + index).fadeIn(500).delay(3000).fadeOut(500, function () {
            fadecount -= 1;
            // すべてfadeoutしたら以下実施
            if (fadecount <= 0) {
                // 完了後のイベント
                $('.navindex' + index).remove();
                // 次があればshowDisp再帰実行
                if (_this.navlist.length > (index + 1)) {
                    _this.showDisp(index + 1);
                }
                else {
                    _this.IsExecute = false;
                }
            }
        });
    };
    Nav.prototype.getMessageHtml = function (nav) {
        return $('<p/>', { 'text': nav.title, 'class': 'navmessagetitle' })[0].outerHTML + $('<p/>', { 'html': nav.message })[0].outerHTML;
    };
    return Nav;
}());
/**
 * querystringに関するクラス
 */
var Query = (function () {
    function Query() {
    }
    /**
     * URL（現在 or 指定の）から、指定の名前の値取得
     * @param key 取得するkey
     * @param url 取得対象のURL。ない場合は現在のURL
     */
    Query.GetParameter = function (key, url) {
        if (!url)
            url = window.location.href;
        key = key.replace(/[\[\]]/g, "\\$&");
        var regex = new RegExp("[?&]" + key + "(=([^&#]*)|&|#|$)"), results = regex.exec(url);
        if (!results)
            return null;
        if (!results[2])
            return '';
        return decodeURIComponent(results[2].replace(/\+/g, " "));
    };
    /**
     * querystring文字列からkeyvalue連想配列取得
     * @param query
     */
    Query.GetArray = function (query) {
        if (!query)
            query = window.location.search;
        if (query.substr(0, 1) == '?')
            query = query.slice(1);
        var hash = query.split('&');
        var result = {};
        for (var i = 0; i < hash.length; i++) {
            var array = hash[i].split('='); //keyと値に分割。
            result[array[0]] = array[1];
        }
        return result;
    };
    /**
     * keyvalue連想配列から、querystring取得※先頭に"?"は付けない
     * @param keyvalue
     */
    Query.GetQueryString = function (keyvalue) {
        var result = '';
        for (var key in keyvalue) {
            if (keyvalue[key] == null || keyvalue[key] == undefined) {
                continue;
            }
            result += (result.length > 0 ? '&' : '');
            result += key + '=' + keyvalue[key];
        }
        return result;
    };
    /**
     * 指定のクエリ文字列から、指定のkey=valueの要素を追加したクエリ文字列を返却
     * ※すでにあれば上書き
     * @param params
     * @param key
     */
    Query.AddParam = function (query, key, value) {
        var array = Query.GetArray(query);
        array[key] = (value != undefined && value != null ? value : '');
        return Query.GetQueryString(array);
    };
    /**
     * 指定のクエリ文字列から、指定のkeyの要素を削除したクエリ文字列を返却
     * @param params
     * @param key
     */
    Query.DeleteParam = function (query, key) {
        var array = Query.GetArray(query);
        if (array[key] != undefined) {
            delete array[key];
        }
        return Query.GetQueryString(array);
    };
    //乗換案内のパラメーター取得
    Query.GetRouteParams = function () {
        Query.RemoveOtherPage();
        var params = {};
        params['mode'] = mode().toString();
        params['frkbn'] = $('#frkbn').val();
        params['frsk'] = $('#frsk').val();
        params['fr'] = encodeURIComponent($('#fr').val());
        params['tokbn'] = $('#tokbn').val();
        params['tosk'] = $('#tosk').val();
        params['to'] = encodeURIComponent($('#to').val());
        if (document.getElementById('nsmode') != null) {
            params['nsmode'] = $('#nsmode').val();
        }
        if (document.getElementsByName('frtokbn') != null) {
            params['frtokbn'] = $("input:radio[name='frtokbn']:checked").val();
        }
        if (document.getElementById('dt_year') != null) {
            params['dt_year'] = $('#dt_year').val();
        }
        if (document.getElementById('dt_month') != null) {
            params['dt_month'] = $('#dt_month').val();
        }
        if (document.getElementById('dt_day') != null) {
            params['dt_day'] = $('#dt_day').val();
        }
        if (document.getElementById('dt_hour') != null) {
            params['dt_hour'] = $('#dt_hour').val();
        }
        if (document.getElementById('dt_minute') != null) {
            params['dt_minute'] = $('#dt_minute').val();
        }
        if (document.getElementById('dt_yyyymm') != null) {
            params['dt_yyyymm'] = $('#dt_yyyymm').val();
        }
        if (document.getElementById('dtkbn') != null) {
            params['dtkbn'] = $('#dtkbn').val();
        }
        if (document.getElementById('sort') != null) {
            params['sort'] = $('#sort').val();
        }
        if (document.getElementById('walkspeed') != null) {
            params['walkspeed'] = $('#walkspeed').val();
        }
        if (document.getElementById('ic') != null) {
            params['ic'] = $('#ic').val();
        }
        if (document.getElementById('idou') != null) {
            params['idou'] = $('#idou').val();
        }
        if (document.getElementById('payrailkbn') != null) {
            params['pr'] = $('#payrailkbn').val();
        }
        if (document.getElementById('frxy') != null) {
            params['frxy'] = $('#frxy').val();
        }
        if (document.getElementById('toxy') != null) {
            params['toxy'] = $('#toxy').val();
        }
        if (document.getElementById('lng') != null) {
            params['lg'] = $('#lng').val();
        }
        if (document.getElementById('frtokbn') != null) {
            params['frtokbn'] = $('#frtokbn').val();
        }
        if (document.getElementById('frtokbn') != null) {
            params['frtokbn'] = $('#frtokbn').val();
        }
        if (document.getElementById('usetrain') != null) {
            params['skbn'] = $('#usetrain').val();
        }
        //custom
        if (document.getElementsByName('kaiseikbn') != null) {
            params['kkbn'] = $("input:radio[name='kaiseikbn']:checked").val();
        }
        return Query.GetQueryString(params);
    };
    Query.RemoveOtherPage = function () {
        // レスポンシブデザイン用にid,nameかぶりが発生し、正常にPOSTされないので、
        // pcであればspを、spであればpcを削除する
        $('.' + (pcsp() == PcSp.pc ? 'sp' : 'pc')).remove();
    };
    return Query;
}());
/**
 * 位置情報
 */
var GeoLocation = (function () {
    function GeoLocation() {
        this.geolocaerr = ['位置情報の利用が許可されていません', '現在地を取得できません', 'タイムアウトしました', '現在地を取得できません'];
        this.positionNames = ['出発地', '到着地', '現在地', '出発地または到着地を選択してください'];
        //ラベル取得
        if (document.getElementById("golocerror") && $('#golocerror').val() != '') {
            this.geolocaerr = $('#golocerror').val().split(':');
        }
        if (document.getElementById("positionNames") && $('#positionNames').val() != '') {
            this.positionNames = $('#positionNames').val().split(':');
        }
    }
    /**
     * 現在座標を取得
     */
    GeoLocation.prototype.GetXy = function () {
        var dfd = $.Deferred();
        // すでに取得済の場合
        if ($('#xy').length > 0) {
            dfd.resolve(true);
            return dfd.promise();
        }
        // 位置情報を使用できる場合
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(function (position) {
                //成功時
                var lat = position.coords.latitude;
                var lng = position.coords.longitude;
                //カスタマイズ
                if (setCurrentLoc(lat, lng)) {
                    var xy = String(lng) + ':' + String(lat);
                    $('body').append($('<input/>', {
                        value: xy,
                        'id': 'xy',
                        'type': 'hidden'
                    }));
                    dfd.resolve(true);
                }
                else {
                    setError(geolocaerr[4]);
                    dfd.reject(false);
                }
                //カスタマイズ end
                //var xy = String(lng) + ':' + String(lat);
                //$('body').append($('<input/>', {
                //    value: xy
                //    , 'id': 'xy'
                //    , 'type': 'hidden'
                //})
                //);
                //dfd.resolve(true);
            }, function (error) {
                //失敗時
                switch (error.code) {
                    case 1:
                        setError(geolocaerr[0]);
                        break;
                    case 2:
                        setError(geolocaerr[1]);
                        break;
                    case 3:
                        setError(geolocaerr[2]);
                        break;
                }
                dfd.reject(false);
            });
        }
        else {
            setError(geolocaerr[3]);
            dfd.reject(false);
        }
        return dfd.promise();
    };
    /**
     * 現在住所を取得
     */
    GeoLocation.prototype.GetAddress = function () {
        var dfd = $.Deferred();
        // 座標未取得の場合エラー
        if ($('#xy').length == 0) {
            dfd.reject(false);
            return dfd.promise();
        }
        var xy = $('#xy').val();
        //住所取得
        $.ajax({
            url: './gaddress',
            cache: false,
            type: 'POST',
            data: { 'qry': encodeURIComponent(xy) },
            timeout: 10000
        }).done(function (data, status, xhr) {
            //成功
            if (data.head.errorcode == '0') {
                var address = data.body.address;
                $('body').append($('<input/>', {
                    value: address,
                    'id': 'address',
                    'type': 'hidden'
                }));
                dfd.resolve(true);
            }
        }).fail(function (xhr, status, error) {
            dfd.reject(false);
            // 失敗
        }).always(function (arg1, status, arg2) {
            // 通信完了
        });
        return dfd.promise();
    };
    return GeoLocation;
}());
var CustomMode = (function () {
    function CustomMode() {
        this.mode = 0;
        this.page = "";
        this.isset_frto = true;
        if ($('#mode_c').length > 0 && $('#mode_c').val().length > 0) {
            var modecustom = $.parseJSON($('#mode_c').val());
            for (var i = 0; i < modecustom.length; i++) {
                if (parseInt($('#mode').val()) == parseInt(modecustom[i]['mode'])) {
                    this.mode = parseInt(modecustom[i]['mode']);
                    this.page = modecustom[i]['page'];
                    this.isset_frto = modecustom[i]['setfrto'] == "1";
                }
            }
        }
    }
    return CustomMode;
}());
if (!String.prototype.format) {
    String.prototype.format = function () {
        var args = arguments;
        return this.replace(/{(\d+)}/g, function (match, number) {
            return typeof args[number] != 'undefined'
                ? args[number]
                : match;
        });
    };
}
///// 値がtrueかどうか
// true、"1"、"True"・・・true返却
// null、undefined、"0"、"False"、その他の値・・・false返却
function parseBool(_obj) {
    if (_obj === null || _obj === "null") {
        return false;
    }
    if (typeof _obj === "undefined") {
        return false;
    }
    if (String(_obj).toLowerCase() == "true") {
        return true;
    }
    if (String(_obj).toLowerCase() == "1") {
        return true;
    }
    return false;
}
function guid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
        s4() + '-' + s4() + s4() + s4();
}
function pcsp() {
    //return $(window).width() < 980 ? PcSp.sp : PcSp.pc;
    return $(window).width() < 768 ? PcSp.sp : PcSp.pc;
}
function pcspSelector() {
    //return $(window).width() < 980 ? PcSp.sp : PcSp.pc;
    return pcsp() == PcSp.pc ? '.pc' : '.sp';
}
function mode() {
    return parseInt($('#mode').val());
}
function is_cmode() {
    if ($('#mode_c').length > 0 && $('#mode_c').val().length > 0) {
        var modecustom = $.parseJSON($('#mode_c').val());
        for (var i = 0; i < modecustom.length; i++) {
            if (parseInt($('#mode').val()) == parseInt(modecustom[i]['mode'])) {
                return true;
            }
        }
    }
    return false;
}
function cmode_data() {
    return new CustomMode();
}
function nsMode() {
    return parseInt($('#nsmode').val());
}
function mapMode() {
    if ($('#mapMode').length > 0) {
        return parseInt($('#mapMode').val());
    }
    return parseInt(Query.GetParameter('mapmode'));
}
/**
 * LIST画面での「戻る」作成
 * @param prms
 */
function pageBack(prms) {
    if (prms === void 0) { prms = null; }
    if (prms == null) {
        prms = $('#params').val();
    }
    //add 2017-02-15 m.ohkoshi ---s
    if (is_cmode()) {
        var cmode = cmode_data();
        location.href = cmode.page + '?' + prms;
        return false;
    }
    //add 2017-02-15 m.ohkoshi ---e
    if (mode() == Mode.diagram) {
        location.href = './diagram?' + prms;
    }
    else if (mode() == Mode.busstate) {
        location.href = './busstate?' + prms;
    }
    else {
        location.href = './route?' + prms;
    }
    return false;
}
/**
 * LIST画面での「戻る」作成
 * ※選択した駅の情報を付与
 * @param prms
 */
function pageBackStation(prms, isfr, st, kbn, stationKbn) {
    var nsmode = parseInt(Query.GetParameter('nsmode'));
    // ここから検索の場合、false(到着)
    if (nsmode == NsMode.quick) {
        isfr = false;
    }
    else if (nsmode == NsMode.direct) {
        // frtokbnが明示的に「1」と指定していればto、
        isfr = !(Query.GetParameter('frtokbn') == '1');
    }
    if (isfr) {
        prms = Query.AddParam(prms, 'fr', encodeURIComponent(st));
        prms = Query.AddParam(prms, 'frkbn', kbn);
        prms = Query.AddParam(prms, 'frsk', stationKbn);
    }
    else {
        prms = Query.AddParam(prms, 'to', encodeURIComponent(st));
        prms = Query.AddParam(prms, 'tokbn', kbn);
        prms = Query.AddParam(prms, 'tosk', stationKbn);
    }
    return pageBack(prms);
}
function setRouteHeader(isfr, st, kbn, skbn, xy) {
    if (xy === void 0) { xy = ''; }
    if (isfr) {
        $('#fr').val(decodeURIComponent(st));
        $('#frkbn').val(kbn);
        $('#frsk').val(skbn);
        $('.frtxt').html(decodeURIComponent(st));
    }
    else {
        $('#to').val(decodeURIComponent(st));
        $('#tokbn').val(kbn);
        $('#tosk').val(skbn);
        $('.totxt').html(decodeURIComponent(st));
    }
    $('html,body').delay(300).animate({
        scrollTop: $('.pageTtl').offset().top
    }, 'slow');
}
/**
 * LIST画面での「戻る」作成
 * ※選択した駅の「出発」「到着」情報を両方付与
 * @param prms
 */
function pageBackStationFrTo(prms, fr, frkbn, frsk, to, tokbn, tosk) {
    if (fr != null) {
        prms = Query.AddParam(prms, 'fr', encodeURIComponent(fr));
        prms = Query.AddParam(prms, 'frkbn', frkbn);
        prms = Query.AddParam(prms, 'frsk', frsk);
    }
    if (to != null) {
        prms = Query.AddParam(prms, 'to', encodeURIComponent(to));
        prms = Query.AddParam(prms, 'tokbn', tokbn);
        prms = Query.AddParam(prms, 'tosk', tosk);
    }
    return pageBack(prms);
}
/**
 * エラーの出力
 * @param msg
 */
function setError(msg) {
    if ($('#jError').length == 0) {
        $('#jErrorField').append($('<div/>', {
            id: 'jError',
            'class': 'jError'
        }));
    }
    if (msg.length > 0) {
        $('#jError').html(msg);
    }
    if (document.getElementById('jError') != null && $('#jError').html().replace(/(^\s+)|(\s+$)/g, "").length > 0) {
        var p = $('.jError').offset().top;
        $('html,body').animate({ scrollTop: p }, 'fast');
    }
}
/**
 * 駅名サジェストイベントの追加
 * @param $input inputフィールド
 * @param afterEventFunc サジェスト実行完了後に、実施するイベント
 */
function AddStationSuggest() {
    var $inputlist = [];
    for (var _i = 0; _i < arguments.length; _i++) {
        $inputlist[_i - 0] = arguments[_i];
    }
    for (var key in $inputlist) {
        var $input = $inputlist[key];
        if ($input.length == 0) {
            return;
        }
        var ekivalues = [];
        var options = {
            // 取得元はajax
            source: function (req, res) {
                if (req.term.length == 0) {
                    res(null);
                    return;
                }
                var data = { 'qry': req.term, 'mode': mode() };
                if ($('[name="kaiseikbn"]:checked').length > 0) {
                    data['knext'] = parseInt($('[name="kaiseikbn"]:checked').val());
                }
                $.ajax({
                    url: "./suggesteki",
                    cache: false,
                    type: 'POST',
                    data: data,
                    dataType: "json",
                    success: function (jsondata) {
                        if (jsondata.head.errorcode == 0 && jsondata.body.ekis.length > 0) {
                            ekivalues = $.map(jsondata.body.ekis, function (a) { return a.value; });
                            console.log(ekivalues);
                            res(jsondata.body.ekis);
                        }
                    }
                });
            },
            search: function (event, ui) {
                var frto = $(event.target).prop('id');
                resetVal(frto);
            },
            select: function (event, ui) {
                $('#' + $(event.target).attr('id')).val('');
                selectEvent(event, ui);
            },
            focus: function (event, ui) {
                //selectEvent(event, ui);//del 2017-01-25
            },
            open: function (event, ui) {
                $('.ui-autocomplete').addClass('historyList suggestList');
                $('.ui-autocomplete').removeClass().addClass('historyList suggestList');
                //$('.ui-menu-item').removeClass('ui-menu-item');
                $('.ui-corner-all').removeClass('ui-corner-all');
            },
            minLength: 0
        };
        //$input.autocomplete(options);
        $input.autocomplete(options).data("ui-autocomplete")._renderItem = function (ul, item) {
            var classname;
            if (item.kubun == StationKbn.Bus) {
                classname = 'icBus';
            }
            else if (item.kubun == StationKbn.Rosen) {
                classname = 'icTrain';
            }
            else {
                classname = 'icLandmark';
            }
            var $li = $("<li>", {
                'class': classname
            }).append('<a href="javascript:void(0);" class="kfr" data-kbn="' + item.kubun + '" data-xy="" data-value="' + item.value + '">' + item.label + '</a>')
                .appendTo(ul);
            return $li;
        };
        $input.on('blur', function () {
            // 要素にないものの場合、空にする
            var frto = $(this).prop('id');
            if (ekivalues.length > 0 && $.inArray($(this).val(), ekivalues) == -1) {
                resetVal(frto);
                return;
            }
        });
        var selectEvent = function (event, ui) {
            var frto = $(event.target).prop('id');
            // 未選択であれば空にする
            if (!ui.item || !ui.item.kubun) {
                resetVal(frto);
                return false;
            }
            $('#' + frto + 'sk').val(ui.item.kubun);
            $('#' + frto + 'img').removeClass().addClass('sticon_' + ui.item.kubun);
            // ランドマークの場合frkbnも追加
            if (ui.item.kubun == StationKbn.Landmark) {
                $('#' + frto + 'kbn').val('6');
            }
            //フォーカス移動 add 2017-01-25
            if (frto == 'fr') {
                $('.history.fr a').focus();
            }
            else {
                $('.history.to a').focus();
            }
            return false;
        };
        var resetVal = function (frto) {
            $('#' + frto + 'sk').val('');
            $('#' + frto + 'img').removeClass();
            $('#' + frto + 'kbn').val('');
        };
    }
}
/// 時刻表表示（方面選択画面、のりば画面で使用）
function viewDiagram(val) {
    var prms = $('#params').val();
    prms = Query.AddParam(prms, 'mode', Mode.diagram);
    // 方面だった場合
    if ($('#diagramgroup').val() == 0) {
        prms = Query.AddParam(prms, 'dgm', encodeURIComponent(val));
        location.href = './diagramdtl?' + prms;
    }
    else {
        prms = Query.AddParam(prms, 'dgmpl', encodeURIComponent(val));
        if ($('#offstationinput').length > 0) {
            prms = Query.AddParam(prms, 'qry', encodeURIComponent($('#offstationinput').val()));
        }
        location.href = './diagrampoledtl?' + prms;
    }
}
/// バス接近情報表示（方面選択画面、のりば画面で使用）
function viewBusState(val) {
    var prms = $('#params').val();
    prms = Query.AddParam(prms, 'mode', Mode.busstate);
    prms = Query.AddParam(prms, 'dgmpl', encodeURIComponent(val));
    location.href = './busstatedtl?' + prms;
}
function getGoogleMapObj(_selector) {
    for (var i = 0; i < $(_selector).length; i++) {
        var map = $(_selector).get(i);
        if (map != null && typeof map.contentWindow != undefined && map.contentWindow.googlemap != undefined) {
            return map.contentWindow.googlemap;
        }
    }
    return null;
}
function onLoadNoriba() {
    // 画像幅を100%に
    $($('#noribaimg')[0].contentWindow.document.body.firstChild).css('width', '100%');
    var cwidth = $('#noribaimg').parent().innerWidth();
    var cheight = $('#noribaimg').parent().innerHeight();
    var imgwidth = $($('#noribaimg')[0].contentWindow.document.body.firstChild).width();
    var imgheight = $($('#noribaimg')[0].contentWindow.document.body.firstChild).height();
    // iframe高さを画像高さにする
    $('#noribaimg').height(imgheight);
    //if (cheight / cwidth < imgheight / imgwidth) {
    //    $($('#noribaimg')[0].contentWindow.document.body.firstChild).css('height', '100%');
    //} else {
    //    $($('#noribaimg')[0].contentWindow.document.body.firstChild).css('width', '100%');
    //}
    $($('#noribaimg')[0].contentWindow.document.body).css('text-align', 'center');
}
// Ajax --------------------------------------------------
/**
 * QRコード画像データをbase64変換した文字列取得
 * @param $div
 * @param _str
 */
function GetQrCodeBase64($div, _str) {
    if ($div.children('img').length > 0) {
        return;
    }
    if (!_str) {
        _str = location.href;
    }
    //QRコード文字列取得
    $.ajax({
        url: './qrcode',
        cache: false,
        type: 'POST',
        data: { 'str': encodeURIComponent(_str) },
        dataType: "json"
    }).done(function (data, status, xhr) {
        //成功
        if (data.head.errorcode == '0') {
            $div.append($('<img/>', {
                'src': 'data:image/png;base64,' + data.body.base64
            }));
        }
    }).fail(function (xhr, status, error) {
        // 失敗
    }).always(function (arg1, status, arg2) {
        // 通信完了
    });
}
/**
 * 検索トップ画面共通項目セット
 * @param _mode
 */
function SetSearchCommon(_mode) {
    // PCSP共通　Myルートアコーディオン開閉
    var toggleBtn = $('.tabArea .toggleBoxTtlBtn');
    toggleBtn.on('click', function () {
        var toggle = $('.toggle');
        $(this).closest('.toggleBox').find(toggle).stop().slideToggle();
        if (toggleBtn.hasClass('open')) {
            $(this).text(openlbl['open']).removeClass('open');
        }
        else {
            $(this).text(openlbl['close']).addClass('open');
        }
    });
    //お知らせ
    var infocls = new Information();
    infocls.SetInformation(openlbl['more'], openlbl['close']);
    //SP時　検索オプション開閉
    var toggleBtn02 = $('.option .toggleBoxTtlBtn');
    toggleBtn02.on('click', function () {
        $(this).closest('.toggleBox').find('.toggle').stop().slideToggle();
        if (toggleBtn02.hasClass('open')) {
            $(this).text(openlbl['close']).removeClass('open');
        }
        else {
            $(this).text(openlbl['open']).addClass('open');
        }
    });
    // 横並びの高さを揃える
    var maxHeight = 0;
    // リスト01
    function listHeight01() {
        var w = $(window).width();
        var list = $('.list01 li');
        var listH = list.height();
        maxHeight = 46;
        list.each(function () {
            if ($(this).height() > maxHeight) {
                maxHeight = $(this).height();
            }
        });
        list.height(maxHeight);
    }
    $(window).on('resize', function () {
        listHeight01();
    });
    listHeight01();
    // タブの縦幅を揃える
    function tabHeight() {
        var w = $(window).width();
        // 同一タブでグループ化しループ
        var $tabgroup = $('.tab').parent();
        for (var i = 0; i < $tabgroup.length; i++) {
            var maxHeight = 36;
            $tabgroup.eq(i).children('.tab').each(function () {
                if (pcsp() == PcSp.pc) {
                    if ($(this).height() > maxHeight) {
                        maxHeight = $(this).height();
                    }
                }
                else {
                    if ($('.tabSearch .tab a').length > 0) {
                        maxHeight = $('.tabSearch .tab a').outerHeight();
                    }
                }
            });
            $tabgroup.eq(i).children('.tab').height(maxHeight);
        }
    }
    $(window).on('resize', function () {
        tabHeight();
    });
    tabHeight();
    function tabWidth() {
        if (pcsp() == PcSp.sp) {
            $('.tabConts .tab').width('32%');
        }
        else {
            $('.tabConts .tab').width('auto');
            var maxwidth = $('.tabConts .tab').width() + 20;
            $.each($('.tabConts .tab'), function () {
                if (maxwidth < $(this).width()) {
                    maxwidth = $(this).width() + 20;
                }
            });
            $('.tabConts .tab').width(maxwidth);
        }
    }
    $(window).on('resize', function () {
        tabWidth();
    });
    tabWidth();
    function tabInfoWidth() {
        if (pcsp() == PcSp.sp) {
            $('.info .tab').width('32%');
        }
        else {
            $('.info .tab').width('auto');
            var maxwidth = $('.info .tab').width() + 20;
            $.each($('.info .tab'), function () {
                if (maxwidth < $(this).width()) {
                    maxwidth = $(this).width() + 20;
                }
            });
            $('.info .tab').width(maxwidth);
        }
    }
    $(window).on('resize', function () {
        tabInfoWidth();
    });
    tabInfoWidth();
    //flash
    function setKohoFlash() {
        $.each($('.historyList'), function () {
            if ($(this).hasClass('koho')) {
                if ($(this).html().trim().length > 0) {
                    $(this).parent().find('.history').addClass('flash3');
                }
                else {
                    $(this).parent().find('.history').addClass('disableBtn');
                    $(this).parent().find('.history').removeClass('mouseon');
                }
            }
        });
    }
    setKohoFlash();
}
function addtHistoryCom(historyval, kbn, kkbn, lang) {
    //localStrage保存 [ { 'lang': '0', 'inputval': [ { 'kbn' : '1' , 'val' : '' , 'xy' : '' , 'skbn' : 'R' } ] } ]
    //json
    var inputobj = [];
    try {
        inputobj = $.parseJSON(localStorage['inputStation' + kbn + kkbn]);
    }
    catch (e) {
        inputobj = [];
    }
    var tmp = [];
    var val = '';
    var vals = [];
    for (var i = 0; i < inputobj.length; i++) {
        if (inputobj[i].lang == lang) {
            vals = inputobj[i].inputval;
            inputobj.splice(i, 1);
        }
    }
    val = historyval.split(';');
    for (var i = 0; i < val.length; i++) {
        if (val[i].split(':').length == 4) {
            tmp.push({ 'kbn': val[i].split(':')[0], 'val': val[i].split(':')[1], 'xy': val[i].split(':')[2].replace(',', ':'), 'skbn': val[i].split(':')[3] });
        }
        else {
            tmp.push({ 'kbn': val[i].split(':')[0], 'val': val[i].split(':')[1], 'xy': '' });
        }
    }
    $.each(vals, function (i, v) {
        var isexists = false;
        for (var j = 0; j < tmp.length; j++) {
            if (v.val == tmp[j].val) {
                isexists = true;
                break;
            }
        }
        if (isexists == false && i <= 8) {
            tmp.push(v);
        }
    });
    try {
        inputobj.push({ 'lang': lang, 'inputval': tmp });
        localStorage['inputStation' + kbn + kkbn] = JSON.stringify(inputobj);
    }
    catch (e) {
    }
}
//カスタマイズ 現在地有効範囲かを判定
function setCurrentLoc(lat, lon) {
    var currentplacelist = $('#currentplacelist').val().split(':');
    //末尾は空白のため削除
    currentplacelist.pop();
    var iseffective = false;
    for (var i in currentplacelist) {
        var coordinate = currentplacelist[i].split(',');
        var minlat = '0';
        var minlon = '0';
        var maxlat = '0';
        var maxlon = '0';
        if (coordinate.length == 4) {
            minlat = coordinate[0];
            minlon = coordinate[1];
            maxlat = coordinate[2];
            maxlon = coordinate[3];
        }
        if ((minlat <= lat && lat <= maxlat) && (minlon <= lon && lon <= maxlon)) {
            iseffective = true;
        }
    }
    return iseffective;
}

/// <reference path="../js/typings/google.analytics/ga.d.ts" />
/// <reference path="../js/typings/jquery/jquery.d.ts" />
(function ($) {
    $(function () {
        if (document.getElementById('analytics_id') == null) {
            return;
        }
        var analytics_id = document.getElementById('analytics_id').value;
        var currdate = new Date();
        (function (i, s, o, g, r, a, m) {
            i['GoogleAnalyticsObject'] = r;
            i[r] = i[r] || function () {
                (i[r].q = i[r].q || []).push(arguments);
            }, i[r].l = 1 * currdate;
            a = s.createElement(o),
                m = s.getElementsByTagName(o)[0];
            a.async = 1;
            a.src = g;
            m.parentNode.insertBefore(a, m);
        })(window, document, 'script', '//www.google-analytics.com/analytics.js', 'ga');
        ga('create', analytics_id, 'auto');
        ga('send', 'pageview');
    });
}(jQuery));
/// <reference path="../js/typings/jquery/jquery.d.ts" />
/// <reference path="../js/typings/jqueryui/jqueryui.d.ts" />
/// <reference path="../js/typings/custom/custom.d.ts" />
/// <reference path="../js/typings/moment/moment.d.ts" />
var MyDiagram = (function () {
    function MyDiagram() {
        this.paramlist = {};
        if ($('#mydiagramparams') != null) {
            var params = Query.GetArray($('#mydiagramparams').val());
            for (var key in params) {
                this.paramlist[key] = decodeURIComponent(params[key]);
            }
        }
    }
    Object.defineProperty(MyDiagram.prototype, "mydiagramid", {
        get: function () { return 'mydiagramArea'; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(MyDiagram.prototype, "$mydiagramdiv", {
        get: function () {
            if ($('#' + this.mydiagramid).length == 0) {
                $('<div/>', {
                    id: this.mydiagramid
                }).appendTo($('body'));
            }
            return $('#' + this.mydiagramid);
        },
        enumerable: true,
        configurable: true
    });
    MyDiagram.prototype.getDiagram = function (targetdivselector) {
        var $targetdiv = $(targetdivselector);
        // My時刻表取得
        var mydiagramjson = localStorage.getItem(buspackage.define.MYDIAGRAM_STORAGE);
        // 何もなければ終了
        if (!mydiagramjson) {
            this.setEmptyResult(targetdivselector);
            return;
        }
        var mydiagramArray = this.getStorageArray();
        var $table = $('<table/>');
        // ヘッダー部分
        $table.append($('<thead/>')
            .append($('<tr/>')
            .append($('<th/>', { 'class': 'md_n_station_company', 'text': this.paramlist['station'] }), $('<th/>', { 'class': 'md_to_name', 'text': this.paramlist['to'] }), $('<th/>', { 'class': 'md_keitou_name', 'text': this.paramlist['keitou'] }), $('<th/>', { 'class': 'md_create_date', 'text': this.paramlist['createdate'] }), $('<th/>', { 'class': 'md_bikou', 'text': this.paramlist['bikou'] }), $('<th/>', { 'class': 'md_buttons', 'text': '' }))));
        var $tbody = $('<tbody/>');
        for (var key in mydiagramArray) {
            var item = mydiagramArray[key];
            var $tr = $('<tr/>', {
                'id': item.guid
            });
            $tr.append($('<td/>', { 'class': 'md_n_station_company', 'text': '' }), $('<td/>', { 'class': 'md_to_name', 'text': '' }), $('<td/>', { 'class': 'md_keitou_name', 'text': '' }), $('<td/>', { 'class': 'md_create_date', 'text': '' }), $('<td/>', { 'class': 'md_bikou', 'text': '' }), $('<td/>', { 'class': 'md_buttons', 'text': '' }));
            $tbody.append($tr);
        }
        $table.append($tbody);
        $targetdiv.append($table);
        this.getDiagramKaisei();
    };
    MyDiagram.prototype.getDiagramKaisei = function () {
        var _this = this;
        var mydiagramjson = localStorage.getItem(buspackage.define.MYDIAGRAM_STORAGE);
        // Ajax
        $.ajax({
            url: "./mydiagram",
            cache: false,
            type: 'POST',
            data: { 'mydiajson': mydiagramjson },
            dataType: "json",
            success: function (jsondata) {
                if (jsondata.head.errorcode == 0 && jsondata.body.mydiagrams.length > 0) {
                    for (var i = 0; i < jsondata.body.mydiagrams.length; i++) {
                        // ajaxからの戻り値を代入
                        var retitem = jsondata.body.mydiagrams[i];
                        if (retitem.isKaisei) {
                            $('.pc').find('#' + retitem.guid).find('.kaisei').first().html(_this.paramlist['message_update']);
                            $('.sp').find('#' + retitem.guid).find('.kaisei').first().html(_this.paramlist['message_update']);
                        }
                    }
                }
            }
        });
    };
    MyDiagram.prototype.getDiagram_top = function (targetdivselector) {
        var $targetdiv = $(targetdivselector);
        var $targetdiv_sp = null;
        // My時刻表取得
        var mydiagramjson = localStorage.getItem(buspackage.define.MYDIAGRAM_STORAGE);
        // 何もなければ終了
        if (!mydiagramjson) {
            var $p = $('<p/>', { 'class': 'mytoolTx', 'text': this.paramlist['nodata'] });
            $targetdiv.append($p);
            return;
        }
        var mydiagramArray = this.getStorageArray();
        for (var key in mydiagramArray) {
            var item = mydiagramArray[key];
            var $div = $('<div/>', { 'class': 'mytoolBox', 'id': item.guid });
            //station
            $div.append($('<p/>', { 'class': 'mytoolTx', 'text': this.paramlist['station'] + ':' })
                .append($('<span/>', { 'class': 'fwb', 'text': item.view_station_name })));
            //noriba
            if (item.pole_no > 0) {
                $div.append($('<p/>', { 'class': 'mytoolTx', 'text': this.paramlist['noriba'] + ':' })
                    .append($('<span/>', { 'class': 'fwb', 'text': item.view_pole_name })));
            }
            //系統番号
            var $keito = $('<span/>', { 'class': 'fwb' });
            var keitos = item.keitou_namelist.split(',');
            for (var j = 0; j < keitos.length; j++) {
                // 同じ名前は追加しない
                if (j > 0) {
                    var checkkeitos = keitos.slice(0, keitos.length - 1);
                    if ($.inArray(keitos[j], checkkeitos) != -1) {
                        continue;
                    }
                }
                var item2 = keitos[j];
                $keito.append(item2);
            }
            $div.append($('<p/>', { 'class': 'mytoolTx', 'text': this.paramlist['keito'] + ':' })
                .append($keito));
            //改正について
            $div.append($('<p/>', { 'text': '', 'class': 'mytoolTx kaisei' }));
            //ボタン
            $div.append($('<ui/>', { 'class': 'btnArea' })
                .append($('<li/>')
                .append($('<a/>', { 'class': 'btn btnSmall mouseon mr10 mydsearch', 'href': 'javascript:void(0);', 'text': this.paramlist['search'] }), $('<a/>', { 'class': 'btn2 btnSmall mouseon myddelete', 'href': 'javascript:void(0);', 'text': this.paramlist['delete'] }))));
            $targetdiv.append($div);
        }
        this.getDiagramKaisei();
        return;
    };
    /**
     * My時刻表情報をタグで取得　※各ページ用
     */
    MyDiagram.prototype.getMyDiagramElements_eachtop = function () {
        var mydiagramArray = this.getStorageArray();
        var ret = '';
        if (mydiagramArray.length > 0) {
            for (var i = 0; i < mydiagramArray.length; i++) {
                var keitos = mydiagramArray[i].keitou_namelist.split(',');
                var keitostr = "";
                for (var j = 0; j < keitos.length; j++) {
                    // 同じ名前は追加しない
                    if (j > 0) {
                        var checkkeitos = keitos.slice(0, keitos.length - 1);
                        if ($.inArray(keitos[j], checkkeitos) != -1) {
                            continue;
                        }
                    }
                    var item2 = keitos[j];
                    keitostr += item2;
                }
                ret += '<div class="toggleList cf" id="' + mydiagramArray[i].guid + '">';
                ret += '<div class="toggleListTx fl">';
                ret += '<dl>';
                ret += '<dt>' + this.paramlist['station'] + '</dt>';
                ret += '<dd>' + mydiagramArray[i].view_station_name + '</dd>';
                ret += '</dl>';
                ret += '<dl>';
                ret += '<dt>' + this.paramlist['noriba'] + '</dt>';
                ret += '<dd>' + mydiagramArray[i].view_pole_name + '</dd>';
                ret += '</dl>';
                ret += '<dl>';
                ret += '<dt>' + this.paramlist['keito'] + '</dt>';
                ret += '<dd>' + keitostr + '</dd>';
                ret += '</dl>';
                ret += '<dl>';
                ret += '<dd class="kaisei"></dd>';
                ret += '</dl>';
                ret += '</div>';
                ret += '<ul class="toggleListBtnArea fr mt20">';
                ret += '<li class="mouseon mr15"><a href="javascript:void(0)" class="mydsearch">' + this.paramlist['search'] + '</a></li>';
                ret += '<li class="deleteBtn mouseon"><a href="javascript:void(0)" class="myddelete">' + this.paramlist['delete'] + '</a></li>';
                ret += '</ul>';
                ret += '</div>';
            }
        }
        else {
            ret += '<div class="toggleList cf" >';
            ret += '<div class="toggleListTx fl">';
            ret += this.paramlist['nodata'];
            ret += '</div>';
            ret += '</div>';
        }
        this.getDiagramKaisei();
        return ret;
    };
    ///**
    // * My時刻表登録ダイアログ表示
    // */
    //public showDiagramDialog(_lng: number, _n_station_company: string, _pole_no: number, _pole_name: string, _knext?: number, _tatewari_keitouname?: string) {
    //    var $dialog = $('<div/>', {
    //        id: this.mydiagramid
    //    }).appendTo($('body'));
    //    //$dialog.append(
    //    //    $('<div/>', { 'class': 'mydiagramContent' })
    //    //    , $('<div/>', { 'class': 'mydiagramButton' })
    //    //);
    //    var buttons: JQueryUI.DialogButtonOptions[] = [];
    //    // My時刻表取得
    //    var mydiagramarray = this.getStorageArray();
    //    // 設定数オーバーの場合
    //    if (mydiagramarray.length >= 10) {
    //        $dialog.html(this.paramlist['error_max']);
    //        buttons = [
    //            {
    //                text: this.paramlist['closebutton']
    //                , click: () => {
    //                    $dialog.dialog('close');
    //                }
    //                , class: 'btnDefault2'
    //            }
    //        ];
    //    }
    //    // 通常時
    //    else {
    //        $dialog.append(
    //            $('<div/>', {
    //                'text': this.paramlist['explain']
    //            })
    //            ,
    //            // タイトル
    //            $('<input/>', {
    //                'type': 'text'
    //                , 'id': 'mydiagramTitle'
    //                , 'placeholder': this.paramlist['placeholder']
    //                , 'value': '{0} {1}'.format(_n_station_company, _pole_name)
    //                , 'class': 'width60p'
    //            })
    //        );
    //        buttons = [
    //            {
    //                text: this.paramlist['savebutton']
    //                , click: () => {
    //                    this.saveDiagramEvent(_lng, _n_station_company, _pole_no, _knext, _tatewari_keitouname);
    //                }
    //                , class: 'btnDefault2'
    //            }
    //            ,
    //            {
    //                text: this.paramlist['cancelbutton']
    //                , click: () => {
    //                    $dialog.dialog('close');
    //                }
    //                , class: 'btnDefault2'
    //            }
    //        ];
    //    }
    //    $dialog.dialog(
    //        {
    //            modal: true
    //            , buttons: buttons
    //            , title: this.paramlist['title']
    //            , close: (event) => {
    //                $dialog.remove();
    //            }
    //            , create: (event, ui) => {
    //                //$("body").css({ overflow: 'hidden' })
    //            }
    //            ,
    //            beforeClose: (event, ui) => {
    //                //$("body").css({ overflow: 'inherit' })
    //            }
    //            , width: 500
    //        });
    //}
    /**
     * My時刻表に登録可能かどうかのチェック
     * @param count 登録予定の件数
     */
    MyDiagram.prototype.CheckDiagramCount = function (count) {
        // My時刻表取得
        var mydiagramarray = this.getStorageArray();
        // 登録後、その件数を超えてしまえばfalse
        if (mydiagramarray.length + count > 10) {
            return false;
        }
        return true;
    };
    /**
     * My時刻表保存イベント
     * @param _lng 登録時言語
     * @param _n_station_company 会社名
     * @param _pole_no 標柱番号
     * @param _tatewari_no 縦割番号
     * @param _keitou_codelist 系統コードリスト
     * @param _kaisei_id 登録時の改正ID
     */
    MyDiagram.prototype.saveDiagramEvent = function (_params) {
        var buttons = [];
        // My時刻表取得
        var mydiagramarray = this.getStorageArray();
        _params.guid = guid();
        _params.title = '';
        _params.create_datetime = new Date();
        // 値を設定
        var mydiagramarray = this.getStorageArray();
        mydiagramarray.push(_params);
        // localstorageに保持
        localStorage.setItem(buspackage.define.MYDIAGRAM_STORAGE, JSON.stringify(mydiagramarray));
        return true;
    };
    /**
     * 表示
     */
    MyDiagram.prototype.moveDiagram = function (guid) {
        var mydiagramArray = this.getStorageArray();
        for (var key in mydiagramArray) {
            var item = mydiagramArray[key];
            //guid合致すれば遷移
            if (item.guid == guid) {
                viewDiagram('{0}:{1}:{2}:{3}'.format(item.n_station_company, item.pole_no, item.tatewari_no, item.keitou_codelist));
                return;
            }
        }
    };
    /**
     * My時刻表削除
     * @param guid
     */
    MyDiagram.prototype.deleteDiagram = function (guid, targetdivselector) {
        var mydiagramArray = this.getStorageArray();
        var newArray = [];
        for (var key in mydiagramArray) {
            var item = mydiagramArray[key];
            //guid合致すれば遷移
            if (item.guid == guid) {
                var $targettr = $('#' + guid);
                $targettr.remove();
            }
            else {
                newArray.push(item);
            }
        }
        // localstorageに保持
        if (newArray.length > 0) {
            localStorage.setItem(buspackage.define.MYDIAGRAM_STORAGE, JSON.stringify(newArray));
        }
        else {
            localStorage.removeItem(buspackage.define.MYDIAGRAM_STORAGE);
            this.setEmptyResult(targetdivselector);
        }
    };
    /**
     * 成功ダイアログ表示
     */
    MyDiagram.prototype.ShowSuccessDialog = function () {
        var _this = this;
        this.$mydiagramdiv.append($('<div/>', {
            'text': this.paramlist['success']
        }));
        var buttons = [
            {
                text: this.paramlist['closebutton'],
                click: function () {
                    $('#' + _this.mydiagramid).dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$mydiagramdiv.dialog({
            modal: true,
            buttons: buttons,
            dialogClass: 'msgdialog',
            title: this.paramlist['title'],
            close: function (event) {
                _this.$mydiagramdiv.remove();
            },
            create: function (event, ui) {
                //$("body").css({ overflow: 'hidden' })
            },
            beforeClose: function (event, ui) {
                //$("body").css({ overflow: 'inherit' })
            },
            width: 500
        });
    };
    /**
     * エラーダイアログ表示
     */
    MyDiagram.prototype.ShowErrorDialog = function () {
        var _this = this;
        var buttons = [
            {
                text: this.paramlist['closebutton'],
                click: function () {
                    _this.$mydiagramdiv.dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$mydiagramdiv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.paramlist['title'],
            close: function (event) {
                _this.$mydiagramdiv.remove();
            },
            create: function (event, ui) {
                //$("body").css({ overflow: 'hidden' })
            },
            beforeClose: function (event, ui) {
                //$("body").css({ overflow: 'inherit' })
            },
            width: 500
        });
    };
    MyDiagram.prototype.getStorageArray = function () {
        // My時刻表取得
        var mydiagramjson = localStorage.getItem(buspackage.define.MYDIAGRAM_STORAGE);
        // 何もなければ終了
        if (!mydiagramjson) {
            return [];
        }
        var mydiagramArray = JSON.parse(mydiagramjson);
        return mydiagramArray;
    };
    /**
     * 0件結果を表示
     */
    MyDiagram.prototype.setEmptyResult = function (targetdivselector) {
        $(targetdivselector).children().remove();
        $(targetdivselector).html(this.paramlist['error_nosave']);
    };
    return MyDiagram;
}());
var MyRoute = (function () {
    function MyRoute() {
        this.labellist = {};
        var labels = Query.GetArray($('#myroutelabels').val());
        for (var key in labels) {
            this.labellist[key] = decodeURIComponent(labels[key]);
        }
        this.geomodel = new GeoLocation();
    }
    Object.defineProperty(MyRoute.prototype, "myrouteid", {
        get: function () { return 'myrouteArea'; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(MyRoute.prototype, "$myroutediv", {
        get: function () {
            if ($('#' + this.myrouteid).length == 0) {
                $('<div/>', {
                    id: this.myrouteid
                }).appendTo($('body'));
            }
            return $('#' + this.myrouteid);
        },
        enumerable: true,
        configurable: true
    });
    /**
     * Myルート保存イベント
     * @param _params
     */
    MyRoute.prototype.saveMyRouteEvent = function (_params) {
        var buttons = [];
        _params.guid = guid();
        _params.create_datetime = new Date();
        //値をセット
        if (!this.checkInput(_params)) {
            return false;
        }
        try {
            //myルート取得
            var myroutearrray = this.getStorageArray();
            myroutearrray.push(_params);
            //localstrage に保持
            localStorage.setItem(buspackage.define.MYROUTE_STORAGE, JSON.stringify(myroutearrray));
        }
        catch (e) {
            this.ShowErrorDialog();
            return false;
        }
        this.ShowSuccessDialog();
        return true;
    };
    MyRoute.prototype.getMyParams = function () {
        var myrouteparam;
        var qryarr = Query.GetArray();
        myrouteparam = {
            guid: guid(),
            create_datetime: new Date(),
            fr: decodeURI(qryarr['fr']),
            to: decodeURI(qryarr['to']),
            frsk: qryarr['frsk'] == null ? '' : qryarr['frsk'],
            tosk: qryarr['tosk'] == null ? '' : qryarr['tosk'],
            lng: $('#lng').val(),
            frxy: qryarr['frxy'] == null ? '' : qryarr['frxy'],
            toxy: qryarr['toxy'] == null ? '' : qryarr['toxy'],
            frHere: qryarr['frkbn'] == String(MenuKbn.here),
            toHere: qryarr['tokbn'] == String(MenuKbn.here)
        };
        // 現在地検索であれば名称変更
        if (myrouteparam.frHere) {
            myrouteparam.fr = this.geomodel.positionNames[2];
        }
        if (myrouteparam.toHere) {
            myrouteparam.to = this.geomodel.positionNames[2];
        }
        return myrouteparam;
    };
    MyRoute.prototype.getStorageArray = function () {
        // Myルート取得
        var myroutejson = localStorage.getItem(buspackage.define.MYROUTE_STORAGE);
        // 何もなければ終了
        if (!myroutejson) {
            return [];
        }
        var myrouteArray = JSON.parse(myroutejson);
        return myrouteArray;
    };
    /**
     * 成功ダイアログ表示
     */
    MyRoute.prototype.ShowSuccessDialog = function () {
        var _this = this;
        this.$myroutediv.append($('<div/>', {
            'text': this.labellist['success']
        }));
        var buttons = [
            {
                text: this.labellist['closebutton'],
                click: function () {
                    $('#' + _this.myrouteid).dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$myroutediv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.labellist['save'],
            close: function (event) {
                _this.$myroutediv.remove();
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            },
            width: 500
        });
    };
    /**
     * エラーダイアログ表示
     */
    MyRoute.prototype.ShowErrorDialog = function () {
        var _this = this;
        var buttons = [
            {
                text: this.labellist['error_nosave'],
                click: function () {
                    _this.$myroutediv.dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$myroutediv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.labellist['save'],
            close: function (event) {
                _this.$myroutediv.remove();
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            },
            width: 500
        });
    };
    /**
    * エラーダイアログ表示
    */
    MyRoute.prototype.ShowMaxErrorDialog = function () {
        var _this = this;
        var buttons = [
            {
                text: this.labellist['error_max'],
                click: function () {
                    _this.$myroutediv.dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$myroutediv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.labellist['save'],
            close: function (event) {
                _this.$myroutediv.remove();
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            },
            width: 500
        });
    };
    /**
     * 入力値チェック
     * @param _params
     */
    MyRoute.prototype.checkInput = function (_params) {
        // Myルート取得
        var myrouteArray = this.getStorageArray();
        var myrouteTemp = [];
        for (var i = 0; i < myrouteArray.length; i++) {
            if (!(myrouteArray[i].fr == _params.fr && myrouteArray[i].to == _params.to && myrouteArray[i].lng == _params.lng)) {
                myrouteTemp.push(myrouteArray[i]);
            }
        }
        localStorage.setItem(buspackage.define.MYROUTE_STORAGE, JSON.stringify(myrouteTemp));
        if (myrouteTemp.length >= 10) {
            this.ShowMaxErrorDialog();
            return false;
        }
        return true;
    };
    /**
     * Myルート情報取得
     */
    MyRoute.prototype.getMyRoute = function () {
        var myrouteArray = this.getStorageArray();
        var myrouteTemp = [];
        var lang = $('#lng').val();
        for (var i = 0; i < myrouteArray.length; i++) {
            if (myrouteArray[i].lng == lang) {
                myrouteTemp.push(myrouteArray[i]);
            }
        }
        return myrouteTemp.reverse();
    };
    /**
     * Myルート情報をタグで取得　※総合トップ用
     */
    MyRoute.prototype.getMyRouteElements_top = function () {
        var myroutes = this.getMyRoute();
        var ret = '';
        if (myroutes.length > 0) {
            for (var i = 0; i < myroutes.length; i++) {
                ret += '<div class="mytoolBox" id="' + myroutes[i].guid + '">';
                ret += '<p class="mytoolTx">' + this.labellist['fromtitle'] + ':<span class="fwb">' + myroutes[i].fr + '</span></p>';
                ret += '<p class="mytoolTx">' + this.labellist['totitle'] + ':<span class="fwb">' + myroutes[i].to + '</span></p>';
                ret += '<ul class="btnArea">';
                ret += '<li><a href="javascript:void(0)" class="btn btnSmall mouseon myroute_nowsearch">' + this.labellist['searchnow'] + '</a></li>';
                ret += '<li><a href="javascript:void(0)" class="btn2 btnSmall mouseon myroute_delete">' + this.labellist['delete'] + '</a></li>';
                ret += '</ul>';
                ret += '</div>';
            }
        }
        else {
            ret = this.labellist['nodata'];
        }
        return ret;
    };
    /**
     * Myルート情報をタグで取得　※各ページ用
     */
    MyRoute.prototype.getMyRouteElements_eachtop = function () {
        var myroutes = this.getMyRoute();
        var ret = '';
        if (myroutes.length > 0) {
            for (var i = 0; i < myroutes.length; i++) {
                ret += '<div class="toggleList cf" id="' + myroutes[i].guid + '">';
                ret += '<div class="toggleListTx fl">';
                ret += '<dl>';
                ret += '<dt>' + this.labellist['fromtitle'] + '</dt>';
                ret += '<dd>' + myroutes[i].fr + '</dd>';
                ret += '</dl>';
                ret += '<dl>';
                ret += '<dt>' + this.labellist['totitle'] + '</dt>';
                ret += '<dd>' + myroutes[i].to + '</dd>';
                ret += '</dl>';
                ret += '</div>';
                ret += '<ul class="toggleListBtnArea fr">';
                ret += '<li class="mouseon"><a href="javascript:void(0)" class="myroute_nowsearch">' + this.labellist['searchnow'] + '</a></li>';
                ret += '<li class="deleteBtn mouseon"><a href="javascript:void(0)" class="myroute_delete">' + this.labellist['delete'] + '</a></li>';
                ret += '</ul>';
                ret += '</div>';
            }
        }
        else {
            ret += '<div class="toggleList cf" >';
            ret += '<div class="toggleListTx fl">';
            ret += this.labellist['nodata'];
            ret += '</div>';
            ret += '</div>';
        }
        return ret;
    };
    /**
     * Myルートを削除
     * @param _guid
     */
    MyRoute.prototype.deleteMyRoute = function (_guid) {
        var myroutes = this.getStorageArray();
        var myrouteTemp = [];
        for (var i = 0; i < myroutes.length; i++) {
            if (myroutes[i].guid != _guid) {
                myrouteTemp.push(myroutes[i]);
            }
        }
        localStorage.setItem(buspackage.define.MYROUTE_STORAGE, JSON.stringify(myrouteTemp));
    };
    /**
     * Myルートの現在時刻で検索パラメータ取得
     * @param _guid
     */
    MyRoute.prototype.getSearchNowParam = function (_guid) {
        var dfd = $.Deferred();
        var myroutes = this.getStorageArray();
        for (var i = 0; i < myroutes.length; i++) {
            if (myroutes[i].guid == _guid) {
                var params = Query.GetArray();
                if (params['dt']) {
                    params['dt'] = '';
                }
                params['fr'] = encodeURIComponent(myroutes[i].fr);
                params['to'] = encodeURIComponent(myroutes[i].to);
                params['frsk'] = myroutes[i].frsk;
                params['tosk'] = myroutes[i].tosk;
                params['frxy'] = myroutes[i].frxy;
                params['toxy'] = myroutes[i].toxy;
                params['frkbn'] = String(MenuKbn.direct);
                params['tokbn'] = String(MenuKbn.direct);
                // 現在地検索の場合、座標取得
                if (myroutes[i].frHere || myroutes[i].toHere) {
                    var isfrhere = myroutes[i].frHere;
                    this.geomodel.GetXy().then(this.geomodel.GetAddress).then(function () {
                        var xy = $('#xy').val();
                        if (isfrhere) {
                            params['frkbn'] = String(MenuKbn.here);
                            params['frxy'] = xy;
                        }
                        else {
                            params['tokbn'] = String(MenuKbn.here);
                            params['toxy'] = xy;
                        }
                        if ($('#address').length > 0) {
                            if (isfrhere) {
                                params['fr'] = encodeURIComponent($('#address').val());
                            }
                            else {
                                params['to'] = encodeURIComponent($('#address').val());
                            }
                        }
                        dfd.resolve(Query.GetQueryString(params));
                    });
                }
                else {
                    dfd.resolve(Query.GetQueryString(params));
                }
            }
        }
        return dfd.promise();
    };
    return MyRoute;
}());
var MyState = (function () {
    function MyState() {
        this.labellist = {};
        var labels = Query.GetArray($('#mystatelabels').val());
        for (var key in labels) {
            this.labellist[key] = decodeURIComponent(labels[key]);
        }
    }
    Object.defineProperty(MyState.prototype, "mystateid", {
        get: function () { return 'mystateArea'; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(MyState.prototype, "$mystatediv", {
        get: function () {
            if ($('#' + this.mystateid).length == 0) {
                $('<div/>', {
                    id: this.mystateid
                }).appendTo($('body'));
            }
            return $('#' + this.mystateid);
        },
        enumerable: true,
        configurable: true
    });
    MyState.prototype.getMyStateData = function () {
        var ret;
        var dgmdl = $('#station_prms').val().split(':');
        //keito
        var keito_code_list = [];
        var keitou_namelist = [];
        for (var i = 0; i < $('.keitou_checkbox').length; i++) {
            var chk = $('.keitou_checkbox')[i];
            if (chk.checked) {
                // 系統コード
                var val = chk.getAttribute('value');
                for (var j = 0; j < val.split(',').length; j++) {
                    keito_code_list.push(val.split(',')[j].split(':')[1]);
                }
                // 系統名
                keitou_namelist.push($(chk).parent().find('.keitou_name').val());
            }
        }
        ret = {
            guid: guid(),
            lng: parseInt($('#lng').val()),
            n_station_name: dgmdl[0],
            pole_no: parseInt(dgmdl[1]),
            view_station_name: dgmdl[2],
            view_pole_name: $('#view_pole_name').val(),
            keito_code_list: keito_code_list.join(','),
            keitou_namelist: keitou_namelist.join(','),
            create_datetime: new Date()
        };
        return ret;
    };
    /**
     * My接近情報保存イベント
     * @param _params
     */
    MyState.prototype.saveMyStateEvent = function (_params) {
        var buttons = [];
        _params.guid = guid();
        _params.create_datetime = new Date();
        //値をセット
        if (!this.checkInput(_params)) {
            return false;
        }
        try {
            //my接近情報取得
            var mystateArray = this.getStorageArray();
            mystateArray.push(_params);
            //localstrage に保持
            localStorage.setItem(buspackage.define.MYBUSSTATE_STORAGE, JSON.stringify(mystateArray));
        }
        catch (e) {
            this.ShowErrorDialog();
            return false;
        }
        this.ShowSuccessDialog();
        return true;
    };
    MyState.prototype.getStorageArray = function () {
        // My接近情報取得
        var mystatejson = localStorage.getItem(buspackage.define.MYBUSSTATE_STORAGE);
        // 何もなければ終了
        if (!mystatejson) {
            return [];
        }
        var mystateArray = JSON.parse(mystatejson);
        return mystateArray;
    };
    /**
     * 成功ダイアログ表示
     */
    MyState.prototype.ShowSuccessDialog = function () {
        var _this = this;
        this.$mystatediv.append($('<div/>', {
            'text': this.labellist['success']
        }));
        var buttons = [
            {
                text: this.labellist['closebutton'],
                click: function () {
                    $('#' + _this.mystateid).dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$mystatediv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.labellist['save'],
            close: function (event) {
                _this.$mystatediv.remove();
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            },
            width: 500
        });
    };
    /**
     * エラーダイアログ表示
     */
    MyState.prototype.ShowErrorDialog = function () {
        var _this = this;
        var buttons = [
            {
                text: this.labellist['error_nosave'],
                click: function () {
                    _this.$mystatediv.dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$mystatediv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.labellist['save'],
            close: function (event) {
                _this.$mystatediv.remove();
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            },
            width: 500
        });
    };
    /**
    * エラーダイアログ表示
    */
    MyState.prototype.ShowMaxErrorDialog = function () {
        var _this = this;
        var buttons = [
            {
                text: this.labellist['error_max'],
                click: function () {
                    _this.$mystatediv.dialog('close');
                },
                class: 'mouseon btn'
            }
        ];
        this.$mystatediv.dialog({
            modal: true,
            dialogClass: 'msgdialog',
            buttons: buttons,
            title: this.labellist['save'],
            close: function (event) {
                _this.$mystatediv.remove();
            },
            create: function (event, ui) {
            },
            beforeClose: function (event, ui) {
            },
            width: 500
        });
    };
    /**
     * 入力値チェック
     * @param _params
     */
    MyState.prototype.checkInput = function (_params) {
        // My接近情報取得
        var mystateArray = this.getStorageArray();
        var mystateTemp = [];
        for (var i = 0; i < mystateArray.length; i++) {
            if (!(mystateArray[i].n_station_name == _params.n_station_name
                && mystateArray[i].pole_no == _params.pole_no
                && mystateArray[i].keito_code_list == _params.keito_code_list
                && mystateArray[i].lng == _params.lng)) {
                mystateTemp.push(mystateArray[i]);
            }
        }
        localStorage.setItem(buspackage.define.MYBUSSTATE_STORAGE, JSON.stringify(mystateTemp));
        if (mystateTemp.length >= 10) {
            this.ShowMaxErrorDialog();
            return false;
        }
        return true;
    };
    /**
     * My接近情報取得
     */
    MyState.prototype.getMyState = function () {
        var mystateArray = this.getStorageArray();
        var mystateTemp = [];
        var lang = $('#lng').val();
        for (var i = 0; i < mystateArray.length; i++) {
            if (mystateArray[i].lng == lang) {
                mystateTemp.push(mystateArray[i]);
            }
        }
        return mystateTemp.reverse();
    };
    /**
     * My接近情報情報をタグで取得　※総合トップ用
     */
    MyState.prototype.getMyStateElements_top = function () {
        var mystates = this.getMyState();
        var ret = '';
        if (mystates.length > 0) {
            for (var i = 0; i < mystates.length; i++) {
                ret += '<div class="mytoolBox" id="' + mystates[i].guid + '">';
                ret += '<p class="mytoolTx">' + this.labellist['stationtitle'] + ':<span class="fwb">' + mystates[i].view_station_name + '</span></p>';
                ret += '<p class="mytoolTx">' + this.labellist['noriba'] + ':<span class="fwb">' + mystates[i].view_pole_name + '</span></p>';
                ret += '<p class="mytoolTx">' + this.labellist['keito'] + ':<span class="fwb">';
                var keitos = mystates[i].keitou_namelist.split(',');
                for (var j = 0; j < keitos.length; j++) {
                    // 同じ名前は追加しない
                    if (j > 0) {
                        var checkkeitos = keitos.slice(0, keitos.length - 1);
                        if ($.inArray(keitos[j], checkkeitos) != -1) {
                            continue;
                        }
                    }
                    var item2 = keitos[j];
                    ret += item2;
                }
                ret += '</span></p>';
                ret += '<ul class="btnArea">';
                ret += '<li><a href="javascript:void(0)" class="btn btnSmall mouseon mystate_search">' + this.labellist['search'] + '</a></li>';
                ret += '<li><a href="javascript:void(0)" class="btn2 btnSmall mouseon mystate_delete">' + this.labellist['delete'] + '</a></li>';
                ret += '</ul>';
                ret += '</div>';
            }
        }
        else {
            ret = this.labellist['nodata'];
        }
        return ret;
    };
    /**
     * My接近情報をタグで取得　※各ページ用
     */
    MyState.prototype.getMyStateElements_eachtop = function () {
        var mystates = this.getMyState();
        var ret = '';
        if (mystates.length > 0) {
            for (var i = 0; i < mystates.length; i++) {
                ret += '<div class="toggleList cf" id="' + mystates[i].guid + '">';
                ret += '<div class="toggleListTx fl">';
                ret += '<dl>';
                ret += '<dt>' + this.labellist['stationtitle'] + '</dt>';
                ret += '<dd>' + mystates[i].view_station_name + '</dd>';
                ret += '</dl>';
                ret += '<dl>';
                ret += '<dt>' + this.labellist['noriba'] + '</dt>';
                ret += '<dd>' + mystates[i].view_pole_name + '</dd>';
                ret += '</dl>';
                ret += '<dl>';
                ret += '<dt>' + this.labellist['keito'] + '</dt>';
                ret += '<dd>';
                var keitos = mystates[i].keitou_namelist.split(',');
                for (var j = 0; j < keitos.length; j++) {
                    // 同じ名前は追加しない
                    if (j > 0) {
                        var checkkeitos = keitos.slice(0, keitos.length - 1);
                        if ($.inArray(keitos[j], checkkeitos) != -1) {
                            continue;
                        }
                    }
                    var item2 = keitos[j];
                    ret += item2;
                }
                ret += '</dd>';
                ret += '</dl>';
                ret += '</div>';
                ret += '<ul class="toggleListBtnArea fr pt20">';
                ret += '<li class="mouseon"><a href="javascript:void(0)" class="mystate_search">' + this.labellist['search'] + '</a></li>';
                ret += '<li class="deleteBtn mouseon"><a href="javascript:void(0)" class="mystate_delete">' + this.labellist['delete'] + '</a></li>';
                ret += '</ul>';
                ret += '</div>';
            }
        }
        else {
            ret += '<div class="toggleList cf" >';
            ret += '<div class="toggleListTx fl">';
            ret += this.labellist['nodata'];
            ret += '</div>';
            ret += '</div>';
        }
        return ret;
    };
    /**
     * My接近情報を削除
     * @param _guid
     */
    MyState.prototype.deleteMyState = function (_guid) {
        var mystates = this.getStorageArray();
        var mystatesTemp = [];
        for (var i = 0; i < mystates.length; i++) {
            if (mystates[i].guid != _guid) {
                mystatesTemp.push(mystates[i]);
            }
        }
        localStorage.setItem(buspackage.define.MYBUSSTATE_STORAGE, JSON.stringify(mystatesTemp));
    };
    /**
     * My接近情報の検索パラメータ取得
     * @param _guid
     */
    MyState.prototype.getSearchParam = function (_guid) {
        var mystates = this.getStorageArray();
        for (var i = 0; i < mystates.length; i++) {
            if (mystates[i].guid == _guid) {
                var params = Query.GetArray();
                var dgmpl = [
                    encodeURIComponent(mystates[i].n_station_name),
                    mystates[i].pole_no,
                    1,
                    encodeURIComponent(mystates[i].keito_code_list)
                ];
                params['dgmpl'] = dgmpl.join(':');
                return Query.GetQueryString(params);
            }
        }
        return '';
    };
    return MyState;
}());

/// <reference path="../js/typings/jqueryui/jqueryui.d.ts" />
/// <reference path="../js/typings/jquery/jquery.d.ts" />
/// <reference path="../js/typings/custom/custom.d.ts" />
var busstatedtl;
/**
 * バス接近情報画面
 */
var BusStateDtl = (function () {
    function BusStateDtl() {
        var _this = this;
        this.btnNames = ['地図を開く', '地図を閉じる'];
        this.openlbl = { 'open': '開く', 'close': '閉じる', 'more': '詳細を見る' };
        /**
        * バス位置情報更新
        */
        this.updateBusState = function () {
            $.ajax({
                url: './busstateupd',
                cache: false,
                type: 'POST',
                data: { 'dgmpl': _this.getCurrentDgmpls(), 'sort4': $('.busstateSortLink.fwb .sort4').val() }
            }).done(function (data, status, xhr) {
                // success
                if (data.head.errorcode == '0') {
                    var buslist = data.body.busstates;
                    var busprmslist = [];
                    var openBuskeyList = []; // 開いているBuskeyのリスト
                    var checkBuskeyList = []; // チェックしているBuskeyのリスト
                    $('.divbusstate').each(function (index, elem) {
                        if ($(elem).hasClass('open')) {
                            openBuskeyList.push($(elem).find('.bus_key').val());
                        }
                    });
                    // チェックしているバスリスト
                    var checklist = _this.getCheckedBusDiv(pcsp() == PcSp.pc ? $('.busstateArea.pc') : $('.busstateArea.sp'));
                    for (var i = 0; i < checklist.length; i++) {
                        checkBuskeyList.push(checklist[i].find('.bus_key').val());
                    }
                    //すでにある内容は削除
                    $('.divbusstate').remove();
                    // 0件時の表示切り替え
                    if (buslist.length == 0) {
                        $('.busstateNoResult').removeClass('notview');
                    }
                    else {
                        $('.busstateNoResult').addClass('notview');
                    }
                    for (var key in buslist) {
                        var bus = buslist[key];
                        var html_pc = $(bus.html);
                        var html_sp = $(bus.html_sp);
                        $('.busstateArea.pc').append(html_pc);
                        $('.busstateArea.sp').append(html_sp);
                        $('#datetimeStr').text(data.body.datetimeStr);
                    }
                    $('.divbusstate').each(function (index, elem) {
                        var key = $(elem).find('.bus_key').val();
                        // 開いていたバスを更新後も開く
                        if (key && $.inArray(key, openBuskeyList) != -1) {
                            $(elem).addClass('open').find('.moreArea').show();
                        }
                        // チェックの変更
                        if (key && $.inArray(key, checkBuskeyList) != -1) {
                            $(elem).find('.busstaetemap').prop('checked', true);
                            busprmslist.push($(elem).find('.busstateprms').val());
                        }
                        else {
                            $(elem).find('.busstaetemap').prop('checked', false);
                        }
                    });
                    _this.busUpdateEvent();
                    var gmapObj = getGoogleMapObj('.map1');
                    if (gmapObj) {
                        gmapObj.ResetBusState(busprmslist);
                    }
                }
            }).fail(function (xhr, status, error) {
            }).always(function (arg1, status, arg2) {
                // 通信完了
            });
        };
        this.busUpdateEvent = function () {
            // 標柱通過アイコンの、幅次第での表示・非表示切り替え
            // ※一瞬moreAreaを表示する、これをやらないと幅がおかしい
            var $ul = $(pcspSelector() + ' .bus_position').parent('ul');
            var $isvisible = $ul.parents('.moreArea').is(':visible');
            if (!$isvisible) {
                $ul.parents('.moreArea').show();
            }
            var bus_positionwidth = $ul.width(); // 標柱通過アイコンの幅
            if (!$isvisible) {
                $ul.parents('.moreArea').hide();
            }
            var w = 0;
            var viewIndex = 0;
            // バス停アイコン・バスアイコンでループ
            for (var i = 0; i < 5; i++) {
                var $item = $(pcspSelector() + ' .bus_position').eq(0).children('.bus_position_item' + i);
                // 要素が無ければ終了
                if ($item.length == 0) {
                    break;
                }
                $item.each(function (index, elem) {
                    w += $(elem).outerWidth(true); // 要素の幅を足す
                });
                // bus_positionの幅内（折りたたまれない）であれば、viewIndexを足す
                if (w <= bus_positionwidth) {
                    viewIndex++;
                }
                else {
                    break;
                }
            }
            // 非表示にする
            for (var i = viewIndex; i < 5; i++) {
                $(' .bus_position').children('.bus_position_item' + i).hide();
            }
            // 詳細を開く
            $('.divbusstate').on('click', function () {
                $(this).toggleClass('open').find('.moreArea').stop().slideToggle();
            });
            // 行押下時
            $('.busstateli').on('click', function (ev) {
                _this.busstateClickEvent($(ev.target));
            });
            // 通過時刻表へ遷移
            $('.btnDiagram').on('click', function (ev) {
                ev.stopPropagation();
                // dgmpl再設定
                var bus_keys = $(ev.target).parents('.divbusstate').find('.bus_key').val().split(':');
                var dgmpls = Query.GetParameter('dgmpl').split(':');
                var hour = Math.floor(parseInt(bus_keys[2]) / 60);
                var minute = parseInt(bus_keys[2]) % 60;
                var newdgmpls = [];
                newdgmpls.push(dgmpls[0]); //停留所名
                newdgmpls.push(dgmpls[1]); //標柱番号
                newdgmpls.push(1); //縦割番号
                newdgmpls.push(bus_keys[3] + bus_keys[1]); //系統コード
                newdgmpls.push(hour); //時間
                newdgmpls.push(minute); //分
                var params = $('#params').val();
                params = Query.AddParam(params, 'dgmpl', encodeURIComponent(newdgmpls.join(':')));
                window.location.href = './linediagram?' + params;
            });
            // バス地図チェックボックス変更時
            // 経路線の変更と、バスアイコン表示変更
            $('.busstaetemap').on('click', function (ev) {
                ev.stopPropagation();
                // チェックしているバスリスト
                var checklist = _this.getCheckedBusDiv(pcsp() == PcSp.pc ? $('.busstateArea.pc') : $('.busstateArea.sp'));
                var busprmslist = [];
                for (var i = 0; i < checklist.length; i++) {
                    busprmslist.push(checklist[i].find('.busstateprms').val());
                }
                var gmapObj = getGoogleMapObj('.map1');
                if (gmapObj) {
                    gmapObj.ResetBusState(busprmslist);
                }
            });
            $('.busstaetemap_label').on('click', function (ev) {
                ev.stopPropagation();
            });
        };
        /**
         * バス位置クリック時イベント
         * @param $targetli クリックしたliのJqueryオブジェクト
         */
        this.busstateClickEvent = function ($targetli) {
            var guid = $targetli.find('.guid').val();
            var gmapObj = getGoogleMapObj('.map1');
            if (gmapObj) {
                gmapObj.FocusBusState(guid);
            }
        };
        if (parseBool($('#viewqrcode').val())) {
            GetQrCodeBase64($('#qrcode'));
        }
        // ボタン名変更
        if ($('#setbtns').val() != null && $('#setbtns').val().split(':').length == 2) {
            this.btnNames = $('#setbtns').val().split(':');
        }
        if ($('#opencloselbl').val() != null) {
            var obj = $.parseJSON(decodeURIComponent($('#opencloselbl').val()));
            this.openlbl = obj;
        }
        // イベント  --------------------------------------------------
        // のりば変更 --------------------------------------------------
        $('#selectPole').on('change', function () {
            var val = $(this).val();
            val = val.split(':').slice(0, 2).join(':');
            var params = $('#params').val();
            params = Query.AddParam(params, 'dgmpl', encodeURIComponent(val));
            location.href = './busstatedtl?' + params;
        });
        //のりばMAPクリック時ページ遷移
        $('#btnNoribaMap').on('click', function () {
            var params = $('#params').val();
            params = Query.AddParam(params, 'mapmode', 2);
            params = Query.AddParam(params, 'mape', encodeURIComponent($('#station_prms').val()));
            location.href = './map?' + params;
        });
        // 表示順変更 --------------------------------------------------
        $('.busstateSortLink').on('click', function () {
            var params = $('#params').val();
            var sort4 = $(this).find('.sort4').val();
            params = Query.AddParam(params, 'sort4', sort4);
            location.href = './busstatedtl?' + params;
        });
        // MAPボタン押下時
        $('.maplink').on('click', function () {
            $(this).toggleClass('open').toggleClass('active');
            var isMapClose = $('.resultMap').css('display') == 'none';
            if (isMapClose) {
                $(this).text(busstatedtl.btnNames[1]);
                busstatedtl.openGooglemapLocation($(this));
            }
            else {
                $(this).text(busstatedtl.btnNames[0]);
            }
            $('.resultMap').stop().slideToggle();
        });
        // SPのみ
        if (pcsp() == PcSp.sp) {
            // 系統一覧の非表示
            $('.typeBoxOuter').hide();
            $('.typeSelectTtl').on('click', function () {
                if ($(this).hasClass('active')) {
                    $(this).removeClass('active');
                    $(this).text(busstatedtl.openlbl['open'] + '▼');
                }
                else {
                    $(this).addClass('active');
                    $(this).text(busstatedtl.openlbl['close'] + '▲');
                }
                $('.typeBoxOuter').stop().slideToggle();
                return false;
            });
        }
        // 系統絞込イベント --------------------------------------------------
        $('.keitou_checkbox').on('change', function (ev) {
            var $this = $(ev.target);
            var tatewarikeitoucodes = $this.val().split(','); // 縦割番号・系統コードのカンマ区切り文字
            for (var i = 0; i < tatewarikeitoucodes.length; i++) {
                var $target = $(" .tatewari_keitoucode[value='" + tatewarikeitoucodes[i] + "']");
                if ($this.prop('checked')) {
                    $target.parent('.divbusstate').removeClass('notview');
                }
                else {
                    $target.parent('.divbusstate').addClass('notview');
                }
            }
            var no = 1;
            $('.pc.divbusstate:not(.notview)').each(function (index, elem) {
                $(elem).find('.bsno').text(no);
                no += 1;
            });
            var no = 1;
            $('.sp.divbusstate:not(.notview)').each(function (index, elem) {
                $(elem).find('.bsno').text(no);
                no += 1;
            });
        });
        var updatesecond = Number($('#autoupdate').val());
        $('.btnReload').on('click', function () {
            // 更新ボタンを押下済みの場合
            if (_this.lastUpdate) {
                var now = new Date();
                // 20秒より前に押下していれば終了
                now.setSeconds(now.getSeconds() - 20);
                if (_this.lastUpdate > now) {
                    return;
                }
            }
            _this.lastUpdate = new Date();
            _this.updateBusState();
        });
        // 時間ごとに自動更新
        if (updatesecond > 0) {
            setInterval(this.updateBusState, updatesecond * 1000);
        }
        // イベントのセット
        this.busUpdateEvent();
        // バス点滅
        setInterval(function () {
            $('.bus_road.bus').toggleClass('view');
        }, 1200);
        //my接近情報
        $('#btnMyBusstate').on('click', function () {
            var myb = new MyState();
            myb.saveMyStateEvent(myb.getMyStateData());
        });
        //add ---s
        if (pcsp() == PcSp.sp && $('#showlist_open').val() == '1') {
            $('.typeSelectTtl').addClass('active');
            $('.typeSelectTtl').text(this.openlbl['close'] + '▲');
            $('.typeBoxOuter').slideToggle();
        }
        //add ---e
    }
    BusStateDtl.prototype.openGooglemapLocation = function (sender) {
        $('.map1').unbind('onload', busstatedtl.map1loaded);
        $('.map1').on('onload', busstatedtl.map1loaded);
        // 標柱情報を取得
        var polelist = $('#mappoleprms').val();
        // バス位置情報を取得
        var busstatelist = [];
        var keitoucode_jougelist = [];
        var syakyokucode = '';
        var checklist = this.getCheckedBusDiv(pcsp() == PcSp.pc ? $('.busstateArea.pc') : $('.busstateArea.sp'));
        for (var i = 0; i < checklist.length; i++) {
            syakyokucode = checklist[i].find('.syakyokucode').val();
            busstatelist.push(checklist[i].find('.busstateprms').val());
            var bus_keys = checklist[i].find('.bus_key').val().split(':');
            var station_prms = $('#station_prms').val().split(':');
            keitoucode_jougelist.push('{0}:{1}:{2}:{3}'.format(bus_keys[0], bus_keys[1], station_prms[3], station_prms[1]));
        }
        // クエリ作成
        var query = 'mode=4&mapmode=4';
        // バス位置情報があれば追加
        if (busstatelist.length > 0) {
            query = Query.AddParam(query, 'buslist', busstatelist.join(','));
        }
        // 系統コードリストがあれば
        if (keitoucode_jougelist.length > 0) {
            query = Query.AddParam(query, 'keitoucode_jougelist', keitoucode_jougelist.join(','));
        }
        else {
            query = Query.AddParam(query, 'polelist', encodeURIComponent(polelist));
        }
        query = Query.AddParam(query, 'syakyokucode', syakyokucode);
        $('.map1')[0].contentDocument.location.replace('./googlemap?' + query);
    };
    BusStateDtl.prototype.map1loaded = function () {
        console.log('map loaded');
    };
    /**
     * 今現在の画面表示（JSで制御したものを含め）に適したパラメータを返却
     * 画面表示している系統のみにする
     */
    BusStateDtl.prototype.getCurrentDgmpls = function () {
        var prms = $('#params').val();
        var keitoucodelist = [];
        $('.keitou_checkbox:checked').each(function (index, elem) {
            var vallist = $(elem).val().split(',');
            {
                for (var i = 0; i < vallist.length; i++) {
                    keitoucodelist.push(vallist[i].split(':')[1]);
                }
            }
        });
        var prms = $('#params').val();
        var dgmpls = Query.GetParameter('dgmpl', prms).split(':');
        return '{0}:{1}::{2}'.format(encodeURIComponent(dgmpls[0]), dgmpls[1], keitoucodelist.join(','));
    };
    /**
     * 選択している地図チェックボックスのdivJquery取得
     */
    BusStateDtl.prototype.getCheckedBusDiv = function ($target) {
        var list = [];
        $target.find('.busstaetemap:checked').each(function (index, elem) {
            list.push($(elem).parents('.divbusstate'));
        });
        return list;
    };
    return BusStateDtl;
}());
$(function () {
    busstatedtl = new BusStateDtl();
});